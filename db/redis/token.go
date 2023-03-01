package redis

import (
	"context"
	"fmt"
	"time"
)

// 保存用户登录时的token,和gpt的accessToken不是一个东西
type TokenCache struct {
	UserID     string `json:"user_id" redis:"user_id"`
	Token      string `json:"token" redis:"token"`
	ExpireTime time.Duration
}

// 把token添加到redis中
func AddTokenCache(ctx context.Context, tokenCache *TokenCache) error {
	//在key中添加token这个字符串是为了和后面保存userInfo的键区分开,因为二者用的都是userID
	key := fmt.Sprintf("token::%s", tokenCache.UserID)
	value := tokenCache.Token
	if err := GetRedisHelper().Set(ctx, key, value, tokenCache.ExpireTime).Err(); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
