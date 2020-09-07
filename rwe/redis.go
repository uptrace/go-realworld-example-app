package rwe

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/go-redis/redisext"
)

var (
	redisRingOnce sync.Once
	redisRing     *redis.Ring
)

func RedisRing() *redis.Ring {
	redisRingOnce.Do(func() {
		opt := Config.RedisCache.Options()
		redisRing = redis.NewRing(opt)

		_ = redisRing.ForEachShard(context.TODO(),
			func(ctx context.Context, shard *redis.Client) error {
				shard.AddHook(redisext.OpenTelemetryHook{})
				return nil
			})
	})
	return redisRing
}

//------------------------------------------------------------------------------

var (
	rateLimiterOnce sync.Once
	rateLimiter     *redis_rate.Limiter
)

func RateLimiter() *redis_rate.Limiter {
	rateLimiterOnce.Do(func() {
		rateLimiter = redis_rate.NewLimiter(RedisRing())
	})
	return rateLimiter
}
