package redishelper

import (
	"context"
	"fmt"
	"log"
	"github.com/go-redis/redis/v8"

    "Event-Delivery-Go/constants"
)

var (
	ctx    = context.Background()
	client *redis.Client
)

// InitializeRedisClient initializes the Redis client.
func InitializeRedisClient() error {
	client = redis.NewClient(&redis.Options{
		Addr:     constants.RedisServerAddress,
		Password: "",
		DB:       0,
	})

	// Test Redis Connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return err
	}
	fmt.Printf("Connected to Redis: %s\n", pong)

	return nil
}

// GetContext returns the context for Redis operations.
func GetContext() context.Context {
	return ctx
}

// PublishToRedisChannel publishes message to Redis channel.
func PublishToRedisChannel(channel, message string) error {
    err := client.Publish(ctx, channel, message).Err()
    if err != nil {
        log.Printf("Failed to publish message to Redis: %v", err)
        return err
    }
    log.Printf("Published message to Redis: %s", message)
    return nil
}

// SubscribeToRedisChannel subscribes to Redis channel.
func SubscribeToRedisChannel(channel string) *redis.PubSub {
	pubsub := client.Subscribe(ctx, channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatalf("Failed to subscribe to Redis channel: %v", err)
	}

	return pubsub
}
