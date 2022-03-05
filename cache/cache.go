package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
}

type CacheImpl struct {
	cli *redis.Client
}

func NewCacheClient(c *redis.Client) *CacheImpl {
	return &CacheImpl{cli: c}
}

// Incr will lock this key and incre
func (c *CacheImpl) Incr(ctx context.Context, key string) error {
	_, err := c.cli.Incr(ctx, key).Result()
	return err
}

// Decr will lock this key and decrement
func (c *CacheImpl) Decr(ctx context.Context, key string) error {
	_, err := c.cli.Decr(ctx, key).Result()
	return err
}

func (c *CacheImpl) Set(ctx context.Context, key string, val int64) error {
	_, err := c.cli.Set(ctx, key, val, 3*time.Minute).Result()
	return err
}

func (c *CacheImpl) Get(ctx context.Context, key string) (string, error) {
	return c.cli.Get(ctx, key).Result()
}

func (c *CacheImpl) LPUSH(ctx context.Context, key string, val ...string) error {
	return c.cli.LPush(ctx, key, val).Err()
}

func (c *CacheImpl) LRANGE(ctx context.Context, key string, start int64, end int64) []string {
	result, _ := c.cli.LRange(ctx, key, start, end).Result()
	return result
}

func (c *CacheImpl) ZADD(ctx context.Context, key string, scoresAndMembers ...interface{}) (int64, error) {
	if len(scoresAndMembers) % 2 != 0 {
		return 0, fmt.Errorf("scores and members is %d, odd", len(scoresAndMembers))
	}

	args := make([]*redis.Z, 0, len(scoresAndMembers)/2)
	for i := 0; i < len(scoresAndMembers); i += 2 {
		score, ok := scoresAndMembers[i].(float64)
		if !ok {
			return 0, fmt.Errorf("score is not float64, %v", scoresAndMembers[i])
		}

		args = append(args, &redis.Z{
			Score:  score,
			Member: scoresAndMembers[i+1],
		})
	}

	return c.cli.ZAdd(ctx, key, args...).Result()
}

func (c *CacheImpl) ZSCAN(ctx context.Context, key string) []string {
	result := make([]string, 0)
	iter := c.cli.ZScan(ctx, key, 0, "", -1).Iterator()
	for  {
		result = append(result, iter.Val())
		if !iter.Next(ctx) {
			break
		}
	}

	return result
}