package helpers

import (
    "log"
	"time"
)


type CustomMessage struct {
    Payload  string
    Attempts int
}

// HandleMessage1 handles messages for Subscriber1.
func HandleMessage1(msg *CustomMessage) bool {
    // Check if the message contains "success1" or "failure1" for failure cases
    if msg.Payload == "failure1" || msg.Payload == "failure12"{
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
func HandleMessage2(msg *CustomMessage) bool {
    // Check if the message contains "good" or "bad" for failure cases
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
