package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Securityevent struct {
	Eventtype string `json:"event_type" binding:"required"`
	SourceIp  string `json:"source_ip" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Timestamp string `json:"timestamp" binding:"required"`
}

func newRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "sentinelai-ingestion",
		})
	})

	router.GET("/readyz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ready": true,
		})
	})

	router.POST("/api/v1/events", func(c *gin.Context) {
		var event Securityevent

		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON or missing required fields",
			})
			return
		}

		if event.Eventtype != "login_failed" &&
			event.Eventtype != "login_success" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unsupported event_type",
			})
			return
		}

		if _, err := time.Parse(time.RFC3339, event.Timestamp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "timestamp must use RFC3339 format",
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message":    "Event accepted",
			"event_type": event.Eventtype,
			"username":   event.Username,
		})
	})

	return router
}
func main() {
	router := newRouter()

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
