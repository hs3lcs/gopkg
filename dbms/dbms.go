package dbms

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

type DBClient struct {
	RW    *sql.DB       // Master
	RO    *sql.DB       // Slave
	Cache *redis.Client // Redis
}

var ctx = context.Background()

func Init(cfg *Config) (*DBClient, error) {
	// Master (RW)
	rwCfg := mysql.Config{
		Net:    "tcp",
		Addr:   cfg.RWDB.HOST,
		User:   cfg.RWDB.USER,
		Passwd: cfg.RWDB.PASS,
	}
	master, err := sql.Open("mysql", rwCfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	// Slave (RO)
	roCfg := mysql.Config{
		Net:    "tcp",
		Addr:   cfg.RODB.HOST,
		User:   cfg.RODB.USER,
		Passwd: cfg.RODB.PASS,
	}
	slave, err := sql.Open("mysql", roCfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	// connection pool
	for _, db := range []*sql.DB{master, slave} {
		db.SetConnMaxLifetime(time.Minute * 1)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		if err := db.Ping(); err != nil {
			return nil, err
		}
	}
	// Redis
	rdb := redis.NewClient(&redis.Options{
		DB:       cfg.REDIS.DB,
		Addr:     cfg.REDIS.HOST,
		Password: cfg.REDIS.PASS,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &DBClient{RW: master, RO: slave, Cache: rdb}, nil
}

func (d *DBClient) SetCache(key string, value any, expire time.Duration) error {
	return d.Cache.Set(ctx, key, value, expire).Err()
}

func (d *DBClient) GetCache(key string) (string, error) {
	return d.Cache.Get(ctx, key).Result()
}

func (d *DBClient) DelCache(key string) error {
	return d.Cache.Del(ctx, key).Err()
}

func (d *DBClient) GetAllCache() (any, error) {
	return d.Cache.Do(ctx, "KEYS", "*").Result()
}
