package wrapperredis

import (
	"context"
	"github.com/hongminhcbg/wrapper-redis/cache"
	"testing"
)

func Test_ZSet(t *testing.T) {
	redisCli := getRedisClint()
	c := cache.NewCacheClient(redisCli)
	ctx := context.Background()
	key := "sset_1"
	i, err := c.ZADD(ctx, key, 0.1, "minh", 2.5, "Hong", 3.7, "nguyen")
	if err != nil {
		t.Error(err)
	}

	t.Log("add key success", i)

	result := c.ZSCAN(ctx, key)
	t.Log("scan key success", result)
}
