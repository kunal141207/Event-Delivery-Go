package helpers

import (
	"log"
	"time"

	"Event-Delivery-Go/redishelper"
	"Event-Delivery-Go/utils"
)

// HandleMessage1 handles messages for Subscriber1.
func HandleMessage1(msg *utils.CustomMessage) bool {
	// Check if the message contains "failure12" or "failure1" for failure cases
	if msg.Payload == "failure1" || msg.Payload == "failure12" {
		time.Sleep(5 * time.Second)
		log.Printf("Subscriber1: Received failure message: %s\n", msg.Payload)
		return false // Failure
	} else {
		time.Sleep(5 * time.Second)
		log.Printf("Subscriber1: Received message: %s\n", msg.Payload)
		return true // Success For Rest of the cases
	}
}

// HandleMessage2 handles messages for Subscriber2.
func HandleMessage2(msg *utils.CustomMessage) bool {
	// Check if the message contains "failure12" or "failure2" for failure cases
	if msg.Payload == "failure2" || msg.Payload == "failure12" {
		time.Sleep(2 * time.Second)
		log.Printf("Subscriber2: Received failure message: %s\n", msg.Payload)
		return false // Failure
	} else {
		time.Sleep(2 * time.Second)
		log.Printf("Subscriber2: Received message: %s\n", msg.Payload)
		return true // Success For Rest of the cases
	}
}

// HandleMessageConsumer handles messages for Subscriber2.
func HandleMessageConsumer(msg *utils.CustomMessage) bool {
	// Check if the message contains "good" or "bad" for failure cases
	if msg.Payload == "failureConsumer" {
		time.Sleep(2 * time.Second)
		log.Printf("Consumer: Received failure message: %s , UserID: %s\n", msg.Payload, msg.UserID)
		return false // Failure
	} else {
		time.Sleep(2 * time.Second)
		log.Printf("Consumer: Received message: %s , UserID: %s\n", msg.Payload, msg.UserID)
		//Publish to redis Pub/Sub for further processing
		err := redishelper.PublishToRedisChannel(utils.RedisChannel, msg.Payload)
		if err != nil {
			log.Printf("Consumer: Received message: %s , UserID: %s\n", msg.Payload, msg.UserID)
			return false
		}
		return true // Success For Rest of the cases
	}
}
