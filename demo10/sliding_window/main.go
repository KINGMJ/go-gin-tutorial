package main

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 滑动窗口限流器结构
type SlidingWindowLimiter struct {
	limit           int           // 窗口请求上限
	windowSize      time.Duration // 滑动窗口时间大小
	smallWindowSize time.Duration // 小窗口时间大小
	numSmallWindows int64         // 小窗口数量
	counters        map[int64]int // 小窗口计数器
	sync.Mutex
}

// 创建滑动窗口限流器
func NewSlidingWindowLimiter(limit int, windowSize, smallWindowSize time.Duration) (*SlidingWindowLimiter, error) {
	// 窗口时间必须能够被小窗口时间整除
	if windowSize%smallWindowSize != 0 {
		return nil, errors.New("window cannot be split by integers")
	}
	return &SlidingWindowLimiter{
		limit:           limit,
		windowSize:      windowSize,
		smallWindowSize: smallWindowSize,
		numSmallWindows: int64(windowSize / smallWindowSize),
		counters:        make(map[int64]int),
	}, nil
}

func (l *SlidingWindowLimiter) TryAcquire() bool {
	l.Lock()
	defer l.Unlock()
	// 获取当前时间
	now := time.Now()

	// 获取当前小窗口值：
	// now.Truncate(l.smallWindow) 是将当前时间 now 截断到最近的小窗口的起始时间。
	// 这意味着它会将当前时间向下舍入到最接近的小窗口的起始时间。
	// 例如，如果当前时间是 2022-01-01 12:34:56.789
	// 而 l.smallWindow 是 5分钟，则 now.Truncate(l.smallWindow) 将会将当前时间截断为 2022-01-01 12:30:00
	currentSmallWindow := now.Truncate(l.smallWindowSize)

	// 获取起始小窗口值
	// 具体来说，currentSmallWindow 是当前小窗口的起始时间
	// l.smallWindowSize 是小窗口的时间大小，l.numSmallWindows 是小窗口的数量。
	// 因此，l.smallWindowSize * time.Duration(l.numSmallWindows-1) 计算了除了当前小窗口外，其余小窗口的时间总和。
	// 然后，通过 Add 方法将当前小窗口的起始时间向前推移这个时间总和，从而得到了窗口的起始时间。
	startSmallWindow := currentSmallWindow.Add(-l.smallWindowSize * time.Duration(l.numSmallWindows-1))

	// 计算当前窗口的请求总数
	var sumCount int

	for smallWindow, count := range l.counters {
		// 清除过期的小窗口计数器，因为滑动窗口每隔一定时间会向右滑动
		// time.Unix(0, smallWindow) 是将一个纳秒级的 Unix 时间戳转换为 time.Time 类型的时间
		if time.Unix(0, smallWindow).Before(startSmallWindow) {
			delete(l.counters, smallWindow)
		} else {
			// 计算当前窗口的请求总数
			sumCount += count
		}
	}
	// 若到达窗口请求上限，请求失败
	if sumCount >= l.limit {
		return false
	}
	// 若没到窗口请求上限，当前小窗口计数器+1，请求成功
	l.counters[currentSmallWindow.UnixNano()]++
	return true
}

func RateLimitMiddleware(rateLimiter *SlidingWindowLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if rateLimiter.TryAcquire() {
			ctx.Next()
		} else {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			ctx.Abort()
		}
	}
}

func main() {
	r := gin.Default()
	// 初始化一个滑动窗口限流器，限制10s最多处理5个请求，每个小窗口为1s
	rateLimiter, err := NewSlidingWindowLimiter(5, 10*time.Second, 1*time.Second)
	if err != nil {
		log.Fatalf("限流器初始化失败：%s", err)
	}
	// 使用自定义中间进行限流
	r.Use(RateLimitMiddleware(rateLimiter))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Run(":8080")
}
