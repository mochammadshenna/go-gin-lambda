package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// RequestLogger middleware using zerolog
func RequestLogger() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		param := gin.LogFormatterParams{
			Request:    c.Request,
			TimeStamp:  time.Now(),
			Latency:    time.Since(start),
			ClientIP:   c.ClientIP(),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
			BodySize:   c.Writer.Size(),
			Keys:       c.Keys,
		}

		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		log.Info().
			Str("method", param.Method).
			Str("path", param.Path).
			Str("ip", param.ClientIP).
			Int("status", param.StatusCode).
			Dur("latency", param.Latency).
			Int("body_size", param.BodySize).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Request processed")
	})
}

// Simple rate limiter
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int // requests per minute
	cleanup  time.Duration
}

type visitor struct {
	requests int
	lastSeen time.Time
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     requestsPerMinute,
		cleanup:  time.Minute * 5,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()
	
	return rl
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(rl.cleanup)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		rl.mu.Lock()
		
		v, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &visitor{
				requests: 1,
				lastSeen: time.Now(),
			}
			rl.mu.Unlock()
			c.Next()
			return
		}

		// Reset counter if more than a minute has passed
		if time.Since(v.lastSeen) > time.Minute {
			v.requests = 1
			v.lastSeen = time.Now()
			rl.mu.Unlock()
			c.Next()
			return
		}

		if v.requests >= rl.rate {
			rl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"code":  http.StatusTooManyRequests,
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		v.requests++
		v.lastSeen = time.Now()
		rl.mu.Unlock()
		c.Next()
	}
}

// ErrorHandler middleware
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Error().
			Interface("error", recovered).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Str("ip", c.ClientIP()).
			Msg("Panic recovered")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"code":    http.StatusInternalServerError,
			"message": "An unexpected error occurred",
		})
	})
}