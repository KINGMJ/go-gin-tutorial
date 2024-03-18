package main

import (
	"sync"
	"time"
)

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	capacity      int       // 桶的容量
	currentTokens int       // 当前令牌数量
	rate          int       // 令牌发放速率/秒
	lastTime      time.Time // 上次发放令牌的时间
	sync.Mutex
}

func NewTokenBucketLimiter(capacity, rate int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity: capacity,
		rate:     rate,
		lastTime: time.Now(),
	}
}

// TryAcquire 尝试获取一个令牌，如果获取成功则返回 true，否则返回 false
func (l *TokenBucketLimiter) TryAcquire() bool {
	l.Lock()
	defer l.Unlock()
	// 获取当前的时间
	now := time.Now()
	// 计算距离上次发放令牌的时间间隔
	interval := now.Sub(l.lastTime)
	// 如果距离上次发放令牌的时间间隔超过1秒，则进行令牌发放操作
	if interval >= time.Second {
		// 计算本次发放令牌的数量，取当前令牌数量和桶的容量的较小值
		l.currentTokens = min(l.capacity, l.currentTokens+int(interval/time.Second)*l.rate)
		// 更新上次发放令牌的时间为当前时间
		l.lastTime = now
	}
	// 如果当前令牌数量为0，则请求失败，否则减少一个令牌并返回成功
	if l.currentTokens == 0 {
		return false
	}
	l.currentTokens--
	return true
}
