package mylog

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// Set example variable
		c.Set("example", "12345")
		c.Set("example2", 1)
		log.Println("---- Before ----")
		// before request
		c.Next()
		log.Println("---- After ----")
		// after request
		latency := time.Since(t)
		log.Print(latency)
		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}
