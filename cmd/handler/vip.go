package handler

import (
	"context"
	"gpt/db/mysql"
	"gpt/internal/global"
	"gpt/pkg/zap"
	"strconv"

	"github.com/gin-gonic/gin"
)

var logger = zap.InitLogger()

// 开通vip
func OpenVIP(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		c.String(400, "user_id不合法")
		return
	}
	//将vip插入到数据库中
	vip := &mysql.VIP{
		UserID:   uint(userId),
		Duration: global.DefaultVIPDuration,
	}
	ctx := context.Background()
	if err := mysql.CreateVIP(ctx, vip); err != nil {
		logger.Error(err.Error())
		c.String(500, "服务器发生错误")
		return
	}
	c.String(200, "开通会员成功")
}
