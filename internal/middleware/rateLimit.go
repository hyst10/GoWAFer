package middleware

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/config"
	"GoWAFer/pkg/utils/api_handler"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

var ipLimitersMutex sync.Mutex
var ipLimiters = make(map[string]RLInterface)

func RateLimitMiddleware(conf *config.Config, repo *repository.IPRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为白名单IP
		if skip, _ := c.Get("isWhiteListIP"); skip == true {
			c.Next()
			return
		}
		// 获取客户端IP
		clientIP := c.ClientIP()
		// 给限速实例map加锁
		ipLimitersMutex.Lock()
		limiter, exists := ipLimiters[clientIP]
		if !exists {
			limiter = NewRateLimiter(conf)
			ipLimiters[clientIP] = limiter
		}
		// 解锁
		ipLimitersMutex.Unlock()
		if !limiter.Allow(clientIP, conf, repo) {
			c.Set("BlockedBy", "CC防护")
			c.Set("BlockReason", "客户端IP访问频率过高")
			api_handler.ForbiddenHandler(c, "访问过于频繁，请稍后再试！")
			c.Abort()
			return
		}
		c.Next()
	}
}

type RLInterface interface {
	Allow(ip string, conf *config.Config, rep *repository.IPRepository) bool
}

// Allow 令牌桶算法
func (rl *TokenBucket) Allow(ip string, conf *config.Config, rep *repository.IPRepository) bool {
	// 检查计数器是否超过最大值
	if rl.Counter > conf.RateLimiter.MaxCounter {
		// 将IP拉入长期黑名单
		expiration := time.Now().AddDate(1, 0, 0)
		current := model.IP{Type: 2, IPAddress: ip, ExpirationAt: expiration}
		rep.Create(&current)
		rl.Counter = 0
		return false
	}
	if rl.Counter > conf.RateLimiter.BanCounter {
		// 将IP封禁一段时间
		expiration := time.Now().Add(time.Duration(conf.RateLimiter.BanDuration) * time.Second)
		current := model.IP{Type: 2, IPAddress: ip, ExpirationAt: expiration}
		rep.Create(&current)
		rl.Counter = 0
		return false
	}
	// 使用互斥锁守护线程
	rl.mu.Lock()
	defer rl.mu.Unlock()
	// 获取当前时间
	now := time.Now()
	// 计算现在距离上次请求过了多少秒
	elapsed := now.Sub(rl.LastTime).Seconds()
	// 将lastTime更新为现在时间，以便下次使用
	rl.LastTime = now
	// 基于经过时间*令牌生成速率来计算需要添加多少令牌
	rl.Tokens += int(elapsed * float64(rl.TokensPerSecond))
	// 当计算结果大于设置的最大令牌数则将当前令牌数设为最大令牌
	if rl.Tokens > rl.MaxTokens {
		rl.Tokens = rl.MaxTokens
	}
	// 当前无令牌，返回false拒绝请求
	if rl.Tokens < 1 {
		rl.Counter++
		return false
	}
	// 执行到这里说明至少存在一个令牌，令牌减1，返回true允许访问
	rl.Tokens--
	return true
}

// Allow 漏桶算法
func (rl *LeakyBucket) Allow(ip string, conf *config.Config, rep *repository.IPRepository) bool {
	// 检查计数器是否超过最大值
	if rl.Counter > conf.RateLimiter.MaxCounter {
		// 将IP拉入长期黑名单
		expiration := time.Now().AddDate(1, 0, 0)
		current := model.IP{Type: 2, IPAddress: ip, ExpirationAt: expiration}
		rep.Create(&current)
		rl.Counter = 0
		return false
	}
	if rl.Counter > conf.RateLimiter.BanCounter {
		// 将IP封禁一段时间
		expiration := time.Now().Add(time.Duration(conf.RateLimiter.BanDuration) * time.Second)
		current := model.IP{Type: 2, IPAddress: ip, ExpirationAt: expiration}
		rep.Create(&current)
		rl.Counter = 0
		return false
	}
	// 使用互斥锁守护线程
	rl.mu.Lock()
	defer rl.mu.Unlock()
	// 获取当前时间
	now := time.Now()
	// 计算当前时间和上次漏水发生的时间差
	elapsed := now.Sub(rl.LastLeakTime)
	// 计算漏水量，从上次调用到现在应该漏掉多少水
	leaks := int(elapsed / rl.LeakInterval)
	// 如果桶里的水比漏掉的水多，从桶里减去漏掉的水，否则将桶里的水设置为0
	if leaks > 0 {
		if rl.Remaining > leaks {
			rl.Remaining -= leaks
		} else {
			rl.Remaining = 0
		}
		// 更新上次漏水发生时间为当前时间
		rl.LastLeakTime = now
	}
	// 检查桶是否有足够的空间来容纳新的水（请求），如果有则增加桶里的水，返回true
	if rl.Remaining < rl.Capacity {
		rl.Remaining++
		return true
	}
	// 桶已装满水，返回false拒绝请求
	rl.Counter++
	return false
}

// Allow 固定窗口算法
func (rl *FixedWindow) Allow(ip string, conf *config.Config, rep *repository.IPRepository) bool {
	// 检查计数器是否超过最大值
	if rl.Counter > conf.RateLimiter.MaxCounter {
		// 将IP拉入长期黑名单
		expiration := time.Now().AddDate(1, 0, 0)
		current := model.IP{Type: 2, IPAddress: ip, ExpirationAt: expiration}
		rep.Create(&current)
		rl.Counter = 0
		return false
	}
	if rl.Counter > conf.RateLimiter.BanCounter {
		// 将IP封禁一段时间
		expiration := time.Now().Add(time.Duration(conf.RateLimiter.BanDuration) * time.Second)
		current := model.IP{Type: 2, IPAddress: ip, ExpirationAt: expiration}
		rep.Create(&current)
		rl.Counter = 0
		return false
	}
	// 使用互斥锁守护线程
	rl.mu.Lock()
	defer rl.mu.Unlock()
	// 获取当前时间
	now := time.Now()
	// 计算距离上次窗口经过的时间，如果大于设定的窗口值，将上次窗口时间设为当前时间，将请求量重置为0
	if now.Sub(rl.WindowStart) >= rl.WindowSize {
		rl.WindowStart = now
		rl.requests = 0
	}
	// 请求量超过了当前窗口范围设定的最大请求量，返回false
	if rl.requests >= rl.MaxRequests {
		return false
	}
	// 请求量加1，返回true
	rl.requests++
	rl.Counter++
	return true
}

// Allow 滑动窗口算法
func (rl *SlidingWindow) Allow(ip string, conf *config.Config, rep *repository.IPRepository) bool {
	// 检查计数器是否超过最大值
	if rl.Counter > conf.RateLimiter.MaxCounter {
		// 将IP拉入长期黑名单
		expiration := time.Now().AddDate(1, 0, 0)
		current := model.IP{Type: 2, IPAddress: ip, ExpirationAt: expiration}
		rep.Create(&current)
		rl.Counter = 0
		return false
	}
	if rl.Counter > conf.RateLimiter.BanCounter {
		// 将IP封禁一段时间
		expiration := time.Now().Add(time.Duration(conf.RateLimiter.BanDuration) * time.Second)
		current := model.IP{Type: 2, IPAddress: ip, ExpirationAt: expiration}
		rep.Create(&current)
		rl.Counter = 0
		return false
	}
	// 使用互斥锁守护线程
	rl.mu.Lock()
	defer rl.mu.Unlock()
	// 获取当前时间
	now := time.Now()
	// 移除早于（当前时间-窗口大小）的时间戳
	var newTimestamps []time.Time
	for _, ts := range rl.timestamps {
		if now.Sub(ts) <= rl.Window {
			newTimestamps = append(newTimestamps, ts)
		}
	}
	// 保存新的时间戳切片
	rl.timestamps = newTimestamps
	// 判断当前请求是否应被允许，如果时间戳队列长度大于最大请求量返回false
	if len(rl.timestamps) >= rl.MaxReq {
		rl.Counter++
		return false
	}
	// 将当前时间戳加入时间戳队列
	rl.timestamps = append(rl.timestamps, now)
	return true
}

// NewRateLimiter 实例化限速器
func NewRateLimiter(conf *config.Config) RLInterface {
	rateMode := conf.RateLimiter.Mode
	switch rateMode {
	case 1:
		return &TokenBucket{
			MaxTokens:       conf.RateLimiter.TokenBucket.MaxToken,
			TokensPerSecond: conf.RateLimiter.TokenBucket.TokenPerSecond,
			Tokens:          conf.RateLimiter.TokenBucket.MaxToken,
			LastTime:        time.Now(),
			Counter:         0,
		}
	case 2:
		return &LeakyBucket{
			Capacity:     conf.RateLimiter.LeakyBucket.Capacity,
			Remaining:    0,
			LeakInterval: time.Second / time.Duration(conf.RateLimiter.LeakyBucket.LeakyPerSecond),
			LastLeakTime: time.Now(),
			Counter:      0,
		}
	case 3:
		return &FixedWindow{
			WindowStart: time.Now(),
			MaxRequests: conf.RateLimiter.FixedWindow.MaxRequest,
			WindowSize:  time.Duration(conf.RateLimiter.FixedWindow.MaxRequest) * time.Second,
			Counter:     0,
		}
	case 4:
		return &SlidingWindow{
			Window:  time.Duration(conf.RateLimiter.SlideWindow.WindowSize) * time.Second,
			MaxReq:  conf.RateLimiter.SlideWindow.MaxRequest,
			Counter: 0,
		}
	default:
		return &TokenBucket{
			MaxTokens:       15,
			TokensPerSecond: 15,
			Tokens:          15,
			LastTime:        time.Now(),
			Counter:         0,
		}
	}
}

// TokenBucket 令牌桶结构体
type TokenBucket struct {
	mu              sync.Mutex
	Tokens          int       // 当前令牌数
	MaxTokens       int       // 最大令牌数
	TokensPerSecond int       // 令牌生成速率
	LastTime        time.Time // 上个令牌生成时间
	Counter         int       // 计数器
}

// LeakyBucket 漏桶结构体
type LeakyBucket struct {
	mu           sync.Mutex
	Capacity     int           // 桶容量
	Remaining    int           // 当前桶里的水
	LeakInterval time.Duration // 漏水速率
	LastLeakTime time.Time     // 上次漏水时间
	Counter      int           // 计数器
}

// FixedWindow 固定窗口结构体
type FixedWindow struct {
	mu          sync.Mutex
	WindowStart time.Time     // 窗口开始时间
	requests    int           // 窗口当前请求数
	MaxRequests int           // 窗口最大请求数
	WindowSize  time.Duration // 窗口大小
	Counter     int           // 计数器
}

// SlidingWindow 滑动窗口结构体
type SlidingWindow struct {
	mu         sync.Mutex
	timestamps []time.Time   // 时间队列
	Window     time.Duration // 窗口大小
	MaxReq     int           // 最大队列数
	Counter    int           // 计数器
}
