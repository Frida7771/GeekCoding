package middlewares

import (
	"GeekCoding/help"
	"GeekCoding/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 时间窗口
	Window time.Duration
	// 允许的最大请求数
	MaxRequests int
	// 限流键前缀
	KeyPrefix string
	// 是否按用户限流（需要从 token 中获取用户信息）
	ByUser bool
	// 是否按 IP 限流
	ByIP bool
}

// RateLimit 限流中间件
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var keys []string

		// 按用户限流
		if config.ByUser {
			// 从 token 中获取用户信息
			auth := c.GetHeader("Authorization")
			if auth != "" {
				userClaims, err := help.AnalyzeToken(auth)
				if err == nil && userClaims != nil {
					// 使用用户 identity 作为限流键
					keys = append(keys, fmt.Sprintf("%s:user:%s", config.KeyPrefix, userClaims.Identity))
				}
			}
		}

		// 按 IP 限流
		if config.ByIP {
			ip := c.ClientIP()
			keys = append(keys, fmt.Sprintf("%s:ip:%s", config.KeyPrefix, ip))
		}

		// 如果没有配置任何限流方式，默认按 IP 限流
		if len(keys) == 0 {
			ip := c.ClientIP()
			keys = append(keys, fmt.Sprintf("%s:ip:%s", config.KeyPrefix, ip))
		}

		// 检查所有限流键
		for _, key := range keys {
			if !checkRateLimit(c, key, config.Window, config.MaxRequests) {
				c.Abort()
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusTooManyRequests,
					"msg":  fmt.Sprintf("请求过于频繁，请稍后再试。限制：%d 次/%v", config.MaxRequests, config.Window),
				})
				return
			}
		}

		c.Next()
	}
}

// checkRateLimit 检查是否超过限流
// 使用滑动窗口算法，通过 Lua 脚本保证原子性
func checkRateLimit(c *gin.Context, key string, window time.Duration, maxRequests int) bool {
	ctx := context.Background()
	now := time.Now()
	windowStart := now.Add(-window).Unix()

	zsetKey := fmt.Sprintf("rate_limit:%s", key)
	member := fmt.Sprintf("%d:%d", now.UnixNano(), now.Nanosecond())
	score := float64(now.Unix())
	expireSeconds := int((window + time.Minute).Seconds())

	// 使用 Lua 脚本保证原子性：清理过期数据、统计、判断、添加、设置过期时间
	luaScript := `
		local key = KEYS[1]
		local window_start = tonumber(ARGV[1])
		local max_requests = tonumber(ARGV[2])
		local score = tonumber(ARGV[3])
		local member = ARGV[4]
		local expire_seconds = tonumber(ARGV[5])
		
		-- 移除窗口外的记录
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start)
		
		-- 获取当前窗口内的请求数
		local count = redis.call('ZCARD', key)
		
		-- 如果超过限制，返回 0（拒绝）
		if count >= max_requests then
			return 0
		end
		
		-- 记录本次请求
		redis.call('ZADD', key, score, member)
		
		-- 设置过期时间
		redis.call('EXPIRE', key, expire_seconds)
		
		-- 返回 1（允许）
		return 1
	`

	result, err := models.RDB.Eval(ctx, luaScript, []string{zsetKey},
		windowStart, maxRequests, score, member, expireSeconds).Result()
	if err != nil {
		// Redis 错误，允许请求通过（降级策略）
		return true
	}

	// Lua 脚本返回 1 表示允许，0 表示拒绝
	return result.(int64) == 1
}

// SubmitRateLimit 提交接口限流（用户级别 + IP 级别）
func SubmitRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Window:      1 * time.Minute, // 1分钟窗口
		MaxRequests: 10,              // 最多10次提交
		KeyPrefix:   "submit",
		ByUser:      true, // 按用户限流
		ByIP:        true, // 按 IP 限流
	})
}

// APIRateLimit 通用 API 限流（IP 级别）
func APIRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Window:      1 * time.Minute, // 1分钟窗口
		MaxRequests: 60,              // 最多60次请求
		KeyPrefix:   "api",
		ByIP:        true, // 按 IP 限流
	})
}

// SendCodeRateLimit 发送验证码限流（IP 级别）
func SendCodeRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Window:      1 * time.Minute, // 1分钟窗口
		MaxRequests: 5,               // 最多5次
		KeyPrefix:   "send_code",
		ByIP:        true, // 按 IP 限流
	})
}

// LoginRateLimit 登录接口限流（IP 级别，防止暴力破解）
func LoginRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Window:      5 * time.Minute, // 5分钟窗口
		MaxRequests: 10,              // 最多10次登录尝试
		KeyPrefix:   "login",
		ByIP:        true, // 按 IP 限流
	})
}

// RegisterRateLimit 注册接口限流（IP 级别，防止批量注册）
func RegisterRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Window:      1 * time.Hour, // 1小时窗口
		MaxRequests: 5,             // 最多5次注册
		KeyPrefix:   "register",
		ByIP:        true, // 按 IP 限流
	})
}

// AdminOperationRateLimit 管理员操作限流（用户级别，防止误操作）
func AdminOperationRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Window:      1 * time.Minute, // 1分钟窗口
		MaxRequests: 20,              // 最多20次操作
		KeyPrefix:   "admin_op",
		ByUser:      true, // 按用户限流
	})
}

// QueryRateLimit 查询接口限流（IP 级别，防止频繁查询）
func QueryRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Window:      1 * time.Minute, // 1分钟窗口
		MaxRequests: 100,             // 最多100次查询
		KeyPrefix:   "query",
		ByIP:        true, // 按 IP 限流
	})
}
