package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 限流器结构体
type RateLimiter struct {
	sync.Mutex               // 互斥锁，用于保护 count 字段
	maxReq     int           // 时间窗口内的最大请求数限制
	interval   time.Duration // 时间窗口间隔
	count      int           // 当前窗口内已处理的请求数
	lastUpdate time.Time     // 上次更新时间
}

func NewRateLimiter(maxReq int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		maxReq:     maxReq,
		interval:   interval,
		lastUpdate: time.Now(),
	}
}

// 2s 100 次
func (r *RateLimiter) Allow() bool {
	r.Lock()
	defer r.Unlock()
	now := time.Now()
	// 当前时间与上次更新时间的时间间隔是否超过了固定窗口的时间间隔
	// 即判断是否已经过了一个固定的时间窗口。
	if now.Sub(r.lastUpdate) >= r.interval {
		r.count = 1
		r.lastUpdate = now
		return true
	}
	// 在时间窗口内是否达到最大限制
	if r.count < r.maxReq {
		r.count++
		return true
	}
	return false
}

func RateLimitMiddleware(rateLimiter *RateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if rateLimiter.Allow() {
			ctx.Next()
		} else {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			ctx.Abort()
		}
	}
}

func main() {
	r := gin.Default()
	// 初始化一个固定窗口限流器，限制10s最多处理5个请求
	rateLimiter := NewRateLimiter(5, 10*time.Second)
	// 使用自定义中间进行限流
	r.Use(RateLimitMiddleware(rateLimiter))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Run(":8080")
}
