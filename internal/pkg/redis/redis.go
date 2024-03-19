package redis_pkg

import (
	"github.com/redis/go-redis/v9"
)

func InitRedis(redisUrl string) (*redis.Client, error) {
	opt := &redis.Options{
		Addr: redisUrl,
	}
	client := redis.NewClient(opt)
	return client, nil
}
