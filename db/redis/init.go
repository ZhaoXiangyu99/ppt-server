package redis

import (
	"context"
	"fmt"
	"gpt/pkg/viper"
	"gpt/pkg/zap"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	config      = viper.Init("db")
	logger      = zap.InitLogger()
	redisOnce   sync.Once
	redisHelper *RedisHelper
)

type RedisHelper struct {
	*redis.Client
}

func GetRedisHelper() *RedisHelper {
	return redisHelper
}

func NewRedisHelper() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Viper.GetString("redis.addr"), config.Viper.GetString("redis.port")),
		Password:     config.Viper.GetString("redis.password"),
		DB:           config.Viper.GetInt("redis.db"),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		//MaxConnAge:   1 * time.Minute,	go-redis v9 已删去
		PoolSize:    10,
		PoolTimeout: 30 * time.Second,
	})

	redisOnce.Do(func() {
		rdh := new(RedisHelper)
		rdh.Client = rdb
		redisHelper = rdh
	})
	return rdb
}

func init() {
	ctx := context.Background()
	rdb := NewRedisHelper()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		logger.Fatalln(err.Error())
		return
	}
	logger.Info("Redis server connection successful!")
	//初始化的时候创建AccessToken池
	AddAccessTokenList(ctx, accessTokens)
}
