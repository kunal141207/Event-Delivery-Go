package subscribers

import (
	"log"
	"time"

    "Event-Delivery-Go/helpers"    
    "Event-Delivery-Go/redishelper" 
    "Event-Delivery-Go/constants"
)


type Subscriber struct {
    Name           string
    messageCh      chan *helpers.CustomMessage
    quitCh         chan struct{}
    handleMessage  func(*helpers.CustomMessage) bool
    maxAttempts    int
    backoffFactor  int
}

// ...

func (s *Subscriber) Start() {
    pubsub := redishelper.SubscribeToRedisChannel(constants.RedisChannel)
	defer pubsub.Close()

	for {
        select {
        case <-s.quitCh:
            // Stop
            return
        default:
            msg, err := pubsub.ReceiveMessage(redishelper.GetContext())
            if err != nil {
                log.Printf("Error receiving message: %v", err)
                continue
            }

            log.Printf("message %s", msg)

            // Process the received message
            customMsg := &helpers.CustomMessage{
                Payload:  msg.Payload,
                Attempts: 0,
            }

            success := s.processMessageWithRetry(customMsg)

            if !success {
                log.Printf("Failed to process message: %s", customMsg.Payload)
            }
        }
    }
}

func (s *Subscriber) processMessageWithRetry(msg *helpers.CustomMessage) bool {
	attempts := 0
	for {
		if s.handleMessage(msg) {
			return true 
		}

		attempts++
		if attempts >= s.maxAttempts {
			return false 
		}

		//backoff with exponential delay
		backoffDuration := time.Duration(1<<uint(msg.Attempts)*s.backoffFactor) * time.Second
		log.Printf("Retry attempt %d in %v seconds...", attempts, backoffDuration.Seconds())
		time.Sleep(backoffDuration)
	}
}

func (s *Subscriber) Stop() {
    close(s.quitCh)
}

func (s *Subscriber) GetName() string {
    return s.Name
}

func NewSubscriber(name string, handleMessage func(*helpers.CustomMessage) bool, maxAttempts, backoffFactor int) *Subscriber {
    return &Subscriber{
        Name:          name,
        messageCh:     make(chan *helpers.CustomMessage),
        quitCh:        make(chan struct{}),
        handleMessage: handleMessage,
        maxAttempts:   maxAttempts,
        backoffFactor: backoffFactor,
    }
}
