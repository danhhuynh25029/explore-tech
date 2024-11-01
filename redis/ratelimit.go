package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"time"
)

var limiter *redis_rate.Limiter

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	limiter = redis_rate.NewLimiter(rdb)
	r := gin.Default()
	r.Use(middleware())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8081")
}

// leaky bucket
func middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := limiter.Allow(c, "project:1231", redis_rate.PerHour(1))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.Header("RateLimit-Remaining", strconv.Itoa(res.Remaining))

		if res.Allowed == 0 {
			seconds := int(res.RetryAfter / time.Second)
			c.Header("RateLimit-RetryAfter", strconv.Itoa(seconds))
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.Next()
	}
}
