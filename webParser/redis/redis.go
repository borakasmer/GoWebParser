package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type redisClient struct {
	c *redis.Client
}

var client = &redisClient{}

var ctx = context.Background()

type ExchangeRedisValue struct{
	Name string
	Value string
	CreatedDate time.Time
}

func GetRedisClient() *redisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}
	client.c = rdb
	return client
}

func (client *redisClient) SetKey(key string, value interface{}, expiration time.Duration) error {
	cacheData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = client.c.Set(ctx, key, cacheData, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *redisClient) GetKey(key string, src interface{}) error {
	val, err := client.c.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}
