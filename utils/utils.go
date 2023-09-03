package utils

const (
    RedisServerAddress = "localhost:6379"
	RedisChannel = "my_channel"
    QueueKey    = "my_queue"
    MaxRetryAttempts   = 3
    BackoffFactor      = 1
)

type CustomMessage struct {
    Payload  string
    Attempts int
    UserID string
}