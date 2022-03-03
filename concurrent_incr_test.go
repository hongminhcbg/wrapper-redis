package wrapperredis

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/hongminhcbg/wrapper-redis/cache"
)

//Test_concurrent this test case to test incr and decr in multiple goroutines, INCR and DECR will be lock key
func Test_concurrent(t *testing.T) {
	c := getRedisClint()
	cacheCli := cache.NewCacheClient(c)
	key := "test_incr_2"
	ctx := context.Background()
	err := cacheCli.Set(ctx, key, 0)
	if err != nil {
		t.Error(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(c *cache.CacheImpl, w *sync.WaitGroup, index int) {
			fmt.Println("start increment key")
			defer wg.Done()
			err := c.Incr(ctx, key)
			if err != nil {
				fmt.Println("start increment key error, firing", err)
				panic(err)
			}

			fmt.Println("start increment key finished ", index)

		}(cacheCli, wg, i)
	}

	for i := 0; i < 100; i++ {
		go func(c *cache.CacheImpl, w *sync.WaitGroup, index int) {
			fmt.Println("start decrement key")
			defer wg.Done()
			err := c.Decr(ctx, key)
			if err != nil {
				fmt.Println("start decrement key error, firing", err)
				panic(err)
			}

			fmt.Println("start decrement key finished ", index)

		}(cacheCli, wg, i)
	}

	wg.Wait()
	result, err := cacheCli.Get(ctx, key)
	if err != nil || result != "0" {
		t.Errorf("error = %v and result is %s", err, result)
	}
}
