package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

const (
	userIDKey = "user_id"
)

type Cache interface {
	InsertChallenge(ctx context.Context, challenge string, userID int64, ttl int) error
	GetChallenge(ctx context.Context, userID int64) (string, error)
	GetUserID(ctx context.Context) (int64, error)
	PresentSusUser(ctx context.Context, userID int64) (bool, error)
	InsertSusUser(ctx context.Context, userID int64, ttl int) error
}

type CacheStruct struct {
	client *redis.Client
}

func NewCacheStruct(client *redis.Client) *CacheStruct {
	return &CacheStruct{
		client: client,
	}
}

func (c *CacheStruct) InsertSusUser(ctx context.Context, userID int64, ttl int) error {
	err := c.client.Set(ctx, fmt.Sprintf("sus_user_%d", userID), "", time.Second*time.Duration(ttl)).Err()
	if err != nil {
		return fmt.Errorf("error saving to redis: %w", err)
	}

	return nil
}

func (c *CacheStruct) PresentSusUser(ctx context.Context, userID int64) (bool, error) {
	_, err := c.client.Get(ctx, fmt.Sprintf("sus_user_%d", userID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, fmt.Errorf("error retieving from redis: %w", err)
	}

	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	return true, nil
}

func (c *CacheStruct) InsertChallenge(ctx context.Context, challenge string, userID int64, ttl int) error {
	err := c.client.Set(ctx, strconv.FormatInt(userID, 10), challenge, time.Second*time.Duration(ttl)).Err()
	if err != nil {
		return fmt.Errorf("error saving to redis: %w", err)
	}

	return nil
}
func (c *CacheStruct) GetChallenge(ctx context.Context, userID int64) (string, error) {
	challenge, err := c.client.Get(ctx, strconv.FormatInt(userID, 10)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("error retieving from redis: %w", err)
	}

	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}

	return challenge, nil
}

func (c *CacheStruct) GetUserID(ctx context.Context) (int64, error) {
	return c.client.Incr(ctx, userIDKey).Result()
}
