package redishelper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"Event-Delivery-Go/utils"
)

var (
	ctx    = context.Background()
	client *redis.Client
)

// InitializeRedisClient initializes the Redis client.
func InitializeRedisClient() error {
	client = redis.NewClient(&redis.Options{
		Addr:     utils.RedisServerAddress,
		Password: utils.RedisPassword,
		DB:       utils.RedisDb,
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

// PublishMessageQueue publishes message to Redis Queue.
func PublishMessageQueue(queueKey string, userID string, payload string) {
	message := utils.CustomMessage{
		UserID:   userID,
		Payload:  payload,
		Attempts: 0,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
	}

	err = pushLeftToQueue(messageJSON, queueKey)
	if err != nil {
		log.Printf("Failed to push message to the queue: %v", err)
	} else {
		log.Printf("Pushed message to the queue: %s", messageJSON)
	}
}

// PublishMessageQueue push left message to Redis Queue.
func pushLeftToQueue(message []byte, QueueKey string) error {
	return client.LPush(ctx, QueueKey, message).Err()
}

// PublishMessageQueue pusj right message to Redis Queue.
func pushRightToQueue(message []byte, QueueKey string) error {
	return client.RPush(ctx, QueueKey, message).Err()
}

// PublishMessageQueue consume message to Redis Queue.
func ConsumerQueue(handleMessage func(msg *utils.CustomMessage) bool, QueueKey string) {
	for {
		result, err := client.BRPop(ctx, 0, QueueKey).Result()
		if err != nil {
			log.Printf("Error popping message from the queue: %v", err)
			continue
		}

		var message *utils.CustomMessage
		err = json.Unmarshal([]byte(result[1]), &message)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		if handleMessage(message) {
			log.Printf("Consumer successful to process message for user %s with payload %s", message.UserID, message.Payload)
			continue
		}

		// Simulate a failure scenario.
		if message.Attempts < utils.MaxRetryAttempts {
			log.Printf("Consumer failed to process message for user %s with payload %s. Retrying...", message.UserID, message.Payload)
			message.Attempts++
			messageJSON, _ := json.Marshal(message)
			pushRightToQueue(messageJSON, QueueKey)
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("MAX RETRIES Consumer processed message for user %s with payload %s", message.UserID, message.Payload)
		}
	}
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
