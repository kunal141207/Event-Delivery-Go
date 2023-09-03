package api

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"Event-Delivery-Go/redishelper"
    "Event-Delivery-Go/constants"
)

// PublishHandler handles publishing data to Redis.
func PublishHandler(c *gin.Context) {
	data := c.PostForm("data")
	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing data parameter"})
		return
	}

	err := redishelper.PublishToRedisChannel(constants.RedisChannel, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data published successfully"})
}
