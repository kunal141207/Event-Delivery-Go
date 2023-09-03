package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"Event-Delivery-Go/api"
	"Event-Delivery-Go/helpers"
	"Event-Delivery-Go/redishelper"
	"Event-Delivery-Go/subscribers"
	"Event-Delivery-Go/utils"
)

func main() {
	// Initialize Redis client
	err := redishelper.InitializeRedisClient()
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	// Initialize and start the subscribers in the background
	subscriber1 := subscribers.NewSubscriber("Subscriber1", helpers.HandleMessage1, utils.MaxRetryAttempts, utils.BackoffFactor)
	subscriber2 := subscribers.NewSubscriber("Subscriber2", helpers.HandleMessage2, utils.MaxRetryAttempts, utils.BackoffFactor)

	go startConsumer()

	go startSubscriber(subscriber1)
	go startSubscriber(subscriber2)

	// Initialize HTTP server
	router := gin.Default()

	// Define API endpoints
	router.POST("/publish", api.PublishHandler)

	// Start HTTP server
	go func() {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Handle signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down...")

	// Stop subscribers and cleanup resources
	stopSubscriber(subscriber1)
	stopSubscriber(subscriber2)
}

func startSubscriber(subscriber *subscribers.Subscriber) {
	log.Printf("Starting %s...\n", subscriber.GetName())
	go subscriber.Start()
}

func stopSubscriber(subscriber *subscribers.Subscriber) {
	log.Printf("Stopping %s...\n", subscriber.GetName())
	subscriber.Stop()
}

func startConsumer() {
	go redishelper.ConsumerQueue(helpers.HandleMessageConsumer, utils.QueueKey)
}
