package wrapperredis

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/hongminhcbg/wrapper-redis/cache"
)

const redisUrl = "redis://localhost:6379"

func getRedisClint() *redis.Client {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	redisCli := redis.NewClient(opts)
	_, err = redisCli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return redisCli
}
func Test_main(t *testing.T) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	redisCli := redis.NewClient(opts)
	result, err := redisCli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("ping redis success", result)
}

func Test_List(t *testing.T) {
	c := getRedisClint()
	cacheCli := cache.NewCacheClient(c)
	key := "test_list_5"
	ctx := context.Background()
	_ = cacheCli.LPUSH(ctx, key, "a")
	_ = cacheCli.LPUSH(ctx, key, "b")
	_ = cacheCli.LPUSH(ctx, key, "cc")
	// now list is c, b, a

	keys := cacheCli.LRANGE(ctx, key, 0, 10)
	if len(keys) != 3 {
		t.Errorf("len keys = %d, keys = %v", len(key), keys)
	}

	t.Log(keys)
}
