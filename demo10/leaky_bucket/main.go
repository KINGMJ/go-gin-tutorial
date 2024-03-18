package main

import (
	"sync"
	"time"
)

type LeakyBucketLimiter struct {
	capacity int       // 桶的容量
	level    int       // 当前水位
	rate     int       // 水流速度/秒
	lastTime time.Time // 上次放水时间
	sync.Mutex
}

func NewLeakyBucketLimiter(capacity, rate int) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		capacity: capacity,
		rate:     rate,
		lastTime: time.Now(),
	}
}

func (l *LeakyBucketLimiter) TryAcquire() bool {
	l.Lock()
	defer l.Unlock()

	// 尝试放水
	now := time.Now()
	// 距离上次放水的时间
	interval := now.Sub(l.lastTime)

	// 判断距离上次放水的时间间隔 interval 是否大于等于 1 秒。
	// 这是为了确保在进行水位更新之前，已经过去了至少 1 秒的时间。
	// 这样做是为了避免过于频繁地更新水位，以免造成不必要的计算开销。
	if interval >= time.Second {
		// 进行水位更新
		// 当前水位 - 距离上次放水的时间(秒)*水流速度
		l.level = max(0, l.level-int(interval/time.Second)*l.rate)
		l.lastTime = now
	}

	// 若到达桶的容量，请求失败
	if l.level >= l.capacity {
		return false
	}
	// 若没有到达桶的容量，当前水位+1，请求成功
	l.level++
	return true
}
