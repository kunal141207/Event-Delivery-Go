package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Event-Delivery-Go/redishelper"
	"Event-Delivery-Go/utils"
)

type InputData struct {
	UserID  string `json:"UserId"`
	Payload string `json:"Payload"`
}

// PublishHandler handles publishing data to Redis.
func PublishHandler(c *gin.Context) {
	var inputData InputData
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redishelper.PublishMessageQueue(utils.QueueKey, inputData.UserID, inputData.Payload)

	c.JSON(http.StatusOK, gin.H{"message": "Data published successfully"})
}
