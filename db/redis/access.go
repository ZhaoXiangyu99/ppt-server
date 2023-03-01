package redis

import (
	"context"
	"errors"
	"gpt/pkg/rand"
)

var (
	accessTokens = []string{
		"",
		"",
	}
)

const listName = "access_token_list"

// 把accessToken全部添加到redis中,采用set来进行存储
func AddAccessTokenList(ctx context.Context, accessTokenList []string) error {
	for _, accessToken := range accessTokenList {
		GetRedisHelper().LPush(ctx, listName, accessToken)
	}
	len, err := GetRedisHelper().LLen(ctx, listName).Result()
	if err != nil || len == 0 {
		return errors.New("Access Token池初始化失败")
	}
	return nil
}

// 从AccessToken池中取出一个来用
func GetAccessToken(ctx context.Context) (string, error) {
	index := rand.GetRand(len(accessTokens))
	accessToken, err := GetRedisHelper().LIndex(ctx, listName, index).Result()
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	return accessToken, nil
}
