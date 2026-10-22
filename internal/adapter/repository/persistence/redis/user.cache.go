package redis

import (

	"context"
	"encoding/json"
	"errors"
	"time"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	goredis "github.com/redis/go-redis/v9"
)

type UserCache struct {
	client *goredis.Client
	ttl   time.Duration
}

func NewUserCache(
	client *goredis.Client,
	ttl   time.Duration,
) *UserCache {
	return &UserCache{client: client, ttl: ttl}
}

const userCacheKeyPrefix = "user:id:"

// Get retrieves a user from the cache by ID. If the user is not found, it returns domain.ErrNotFound.
func (c *UserCache) Get(
	ctx context.Context, id string,
) (*domain.User, error) {

	// Attempt to get the user data from Redis
	data, err := c.client.Get(ctx, userCacheKeyPrefix + id).Bytes()
	if errors.Is(err, goredis.Nil) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a domain.User object
	var u domain.User
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

// Set stores a user in the cache with a specified TTL (time-to-live).
func(c *UserCache) Set(ctx context.Context, u *domain.User) error{
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	uuid := u.ID.String()

	return c.client.Set(ctx, userCacheKeyPrefix + uuid, data, c.ttl).Err()
}

// Delete removes a user from the cache by ID.
func(c *UserCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, userCacheKeyPrefix + id).Err()
}