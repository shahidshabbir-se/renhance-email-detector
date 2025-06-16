package db

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/shahidshabbir-se/renhance-email-detector/pkg/utils"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func InitRedis(ctx context.Context, log *logrus.Logger) error {
	addr := utils.GetEnv("REDIS_ADDR", "localhost:6379")
	password := utils.GetEnv("REDIS_PASSWORD", "")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	log.Info("Redis connected")
	return nil
}
