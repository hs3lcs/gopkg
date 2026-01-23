package dbms

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

var ctx = context.Background()

func rdbConn(db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     Config.REDIS_HOST,
		Password: Config.REDIS_PASS,
		DB:       db,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})
	return rdb
}

func CacheSet(key, value string, expire time.Duration) error {
	rdb := rdbConn(0)
	err := rdb.Set(ctx, key, value, expire).Err()
	return err
}

func CacheGet(key string) (string, error) {
	rdb := rdbConn(0)
	res, err := rdb.Get(ctx, key).Result()
	return res, err
}

func CacheDel(key string) error {
	rdb := rdbConn(0)
	err := rdb.Del(ctx, key).Err()
	return err
}

func CacheGetAll() (any, error) {
	rdb := rdbConn(0)
	res, err := rdb.Do(ctx, "KEYS", "*").Result()
	return res, err
}
