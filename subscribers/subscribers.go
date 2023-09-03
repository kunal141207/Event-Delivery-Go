package subscribers

import (
	"log"
	"time"

	"Event-Delivery-Go/redishelper"
	"Event-Delivery-Go/utils"
)

type Subscriber struct {
	Name          string
	messageCh     chan *utils.CustomMessage
	quitCh        chan struct{}
	handleMessage func(*utils.CustomMessage) bool
	maxAttempts   int
	backoffFactor int
}

// ...

func (s *Subscriber) Start() {
	pubsub := redishelper.SubscribeToRedisChannel(utils.RedisChannel)
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
			customMsg := &utils.CustomMessage{
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

func (s *Subscriber) processMessageWithRetry(msg *utils.CustomMessage) bool {
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

func NewSubscriber(name string, handleMessage func(*utils.CustomMessage) bool, maxAttempts, backoffFactor int) *Subscriber {
	return &Subscriber{
		Name:          name,
		messageCh:     make(chan *utils.CustomMessage),
		quitCh:        make(chan struct{}),
		handleMessage: handleMessage,
		maxAttempts:   maxAttempts,
		backoffFactor: backoffFactor,
	}
}
