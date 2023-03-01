package handler

import (
	"context"
	"encoding/json"
	"gpt/internal/request"
	"gpt/internal/response"
	"gpt/pkg/jwt"
	"strconv"
	"time"

	"gpt/db/mysql"
	"gpt/db/redis"

	"github.com/gin-gonic/gin"
)

// 创建新的对话
func NewConversation(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		c.String(400, "user_id不合法")
		return
	}
	//数据库中查询vip是否存在
	ctx := context.Background()
	vip, err := mysql.GetVIPByUserID(ctx, userId)
	if err != nil {
		logger.Error(err.Error())
		c.String(500, "用户不是vip")
		return
	}

	//生成token
	expireTime := time.Now().AddDate(0, 0, int(vip.Duration))
	Jwt := jwt.NewJWT([]byte("signKey"))
	claims := jwt.CustomClaims{UserID: userId}
	claims.ExpiresAt = expireTime.Unix()
	token, err := Jwt.CreateToken(claims)
	if err != nil {
		logger.Error(err.Error())
		c.String(500, "服务器发生错误")
		return
	}

	//把token保存到redis中
	tokenCache := &redis.TokenCache{
		UserID:     userIdStr,
		Token:      token,
		ExpireTime: time.Hour * 24 * time.Duration(vip.Duration),
	}
	err = redis.AddTokenCache(ctx, tokenCache)
	if err != nil {
		logger.Error(err.Error())
		c.String(500, "服务器发生错误")
		return
	}

	//把token返回给用户
	res := response.ConversationResponse{Token: token}
	c.JSON(200, res)
}

// 开始对话(已经创建过对话)
func StartConversation(c *gin.Context) {
	//从header中获取token并校验token
	token := c.GetHeader("Authorization")
	Jwt := jwt.NewJWT([]byte("signKey"))
	claims, err := Jwt.ParseToken(token)
	if err != nil {
		logger.Error(err.Error())
		c.String(403, "token解析错误")
		return
	}
	userId := claims.UserID
	ctx := context.Background()
	//获取前端的参数,绑定json
	var param *request.ConversationParam
	if err := c.BindJSON(param); err != nil {
		logger.Error(err.Error())
		c.String(400, "参数个数错误")
		return
	}
	//获取用户的accessToken
	var accessToken string
	//如果用户的信息在redis中没有保存,那么就将用户的信息保存到redis中,并且从池子里取出一个accessToken
	if userCache, err := redis.GetUserCache(ctx, strconv.Itoa(int(userId))); err != nil {
		//用户信息不存在
		accessToken, err = redis.GetAccessToken(ctx)
		if err != nil {
			logger.Error("获取accessToken失败")
			c.String(500, "服务器发生错误")
			return
		}
		//把用户信息保存到redis中,下次对话时直接从redis中读取
		userCache := &redis.UserCache{
			UserID:          strconv.Itoa(int(userId)),
			ConversationId:  param.Options.ConversationId,
			ParentMessageId: param.Options.ParentMessageId,
			AccessToken:     accessToken,
		}
		if err := redis.AddUserCache(ctx, userCache); err != nil {
			logger.Error("添加用户信息错误")
			c.String(500, "服务器发生错误")
			return
		}
	} else {
		//用户信息在redis中已存在,直接从redis中取出来用(主要是AccessToken)
		//将param重新赋值
		param.Options = request.Option{
			ConversationId:  userCache.ConversationId,
			ParentMessageId: userCache.ParentMessageId,
		}
		//将accessToken赋值为redis中保存的
		accessToken = userCache.AccessToken
	}

	//将参数给传给python || node server
	jsonParam, _ := json.Marshal(param)
	body, err := request.SendPost(jsonParam, accessToken)
	if err != nil {
		c.String(500, "服务器发生错误")
		return
	}
	var resp response.GPT
	if err = json.Unmarshal(body, &resp); err != nil {
		logger.Error(err.Error())
		c.String(500, "服务器发生错误")
		return
	}

	c.JSON(200, resp)
}
