package backends

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lingzhi227/go-scheduler/pkg/memory"
	"github.com/redis/go-redis/v9"
)

// Redis implements memory.Backend using Redis for distributed storage.
type Redis struct {
	client *redis.Client
	prefix string
}

// NewRedis creates a Redis-backed memory store.
func NewRedis(addr, password string, db int) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
		prefix: "agentsched:mem:",
	}
}

// NewRedisFromClient creates a Redis backend from an existing client.
func NewRedisFromClient(client *redis.Client, prefix string) *Redis {
	if prefix == "" {
		prefix = "agentsched:mem:"
	}
	return &Redis{client: client, prefix: prefix}
}

func (r *Redis) redisKey(scope memory.Scope, scopeID, key string) string {
	return fmt.Sprintf("%s%s:%s:%s", r.prefix, scope, scopeID, key)
}

func (r *Redis) setKey(scope memory.Scope, scopeID string) string {
	return fmt.Sprintf("%s%s:%s:__keys__", r.prefix, scope, scopeID)
}

func (r *Redis) Set(scope memory.Scope, scopeID, key string, value any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}

	pipe := r.client.Pipeline()
	pipe.Set(ctx, r.redisKey(scope, scopeID, key), data, 0)
	pipe.SAdd(ctx, r.setKey(scope, scopeID), key)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *Redis) Get(scope memory.Scope, scopeID, key string) (any, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := r.client.Get(ctx, r.redisKey(scope, scopeID, key)).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, false, err
	}
	return value, true, nil
}

func (r *Redis) Delete(scope memory.Scope, scopeID, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipe := r.client.Pipeline()
	pipe.Del(ctx, r.redisKey(scope, scopeID, key))
	pipe.SRem(ctx, r.setKey(scope, scopeID), key)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *Redis) List(scope memory.Scope, scopeID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.SMembers(ctx, r.setKey(scope, scopeID)).Result()
}

func (r *Redis) SearchVector(_ memory.Scope, _ string, _ []float64, _ int) ([]memory.VectorResult, error) {
	// Redis vector search requires RediSearch module.
	// For now, return empty; production would use FT.SEARCH.
	return nil, nil
}

// Ping checks Redis connectivity.
func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the Redis client.
func (r *Redis) Close() error {
	return r.client.Close()
}
