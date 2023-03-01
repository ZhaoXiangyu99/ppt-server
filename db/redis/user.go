package redis

import (
	"context"
	"fmt"
	"time"
)

// 保存用户用于向GPT发送请求的AccessToken和ConversationId
type UserCache struct {
	UserID          string `json:"user_id" redis:"user_id"`
	AccessToken     string `json:"access_token" redis:"access_token"`
	ConversationId  string `json:"conversationId" redis:"conversation_id"`
	ParentMessageId string `json:"parentMessageId" redis:"parrent_message_id"`
	expireTime      time.Duration
}

// 添加用户的重要信息到redis
func AddUserCache(ctx context.Context, userCache *UserCache) error {
	key := fmt.Sprintf("info::%s", userCache.UserID)
	value := make(map[string]string)
	value["access_token"] = userCache.AccessToken
	value["conversation_id"] = userCache.ConversationId
	value["parrent_message_id"] = userCache.ParentMessageId
	if err := GetRedisHelper().HSet(ctx, key, value).Err(); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// 删除redis中用户的信息
func DeleteUserCache(ctx context.Context, userId string) error {
	key := fmt.Sprintf("info::%s", userId)
	if err := GetRedisHelper().Del(ctx, key).Err(); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// 查询用户信息是否存在
func GetUserCache(ctx context.Context, userId string) (*UserCache, error) {
	key := fmt.Sprintf("info::%s", userId)
	result, err := GetRedisHelper().HGetAll(ctx, key).Result()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	uc := &UserCache{
		UserID:          userId,
		AccessToken:     result["access_token"],
		ConversationId:  result["conversation_id"],
		ParentMessageId: result["parrent_message_id"],
	}
	return uc, nil
}
