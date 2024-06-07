package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		latency := time.Since(t)
		clienIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		userUrgent := ctx.Request.UserAgent()
		path := ctx.Request.URL.Path

		log.Printf("[LOG] %s -[%v] \"%s %s %d %v \"%s\"\n", clienIP, t, method, path, statusCode, latency, userUrgent)
	}
}
