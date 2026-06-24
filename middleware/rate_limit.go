package middleware

import (
	"fmt"
	"invitation-app/config"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	libredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// Cache limiter supaya tidak dibuat ulang tiap request
var (
	limiters   = make(map[string]*limiter.Limiter)
	limitersMu sync.RWMutex
)

// Inisialisasi store — otomatis pilih Redis atau Memory
func newStore() limiter.Store {
	if config.RedisClient != nil {
		store, err := libredis.NewStoreWithOptions(
			config.RedisClient,
			limiter.StoreOptions{
				Prefix:          "invitation_rl",  // prefix key di Redis
				MaxRetry:        3,                // retry jika Redis timeout
				CleanUpInterval: 5 * time.Minute, // bersihkan key expired
			},
		)
		if err != nil {
			// Fallback ke memory jika Redis gagal
			fmt.Println("⚠️  Redis store gagal, fallback ke memory store:", err)
			return memory.NewStore()
		}
		return store
	}

	// Redis tidak tersedia → pakai memory
	fmt.Println("⚠️  Redis tidak tersedia, pakai memory store")
	return memory.NewStore()
}

// Helper — ambil atau buat limiter berdasarkan nama unik
func getLimiter(name string, period time.Duration, limit int64) *limiter.Limiter {
	limitersMu.RLock()
	if l, ok := limiters[name]; ok {
		limitersMu.RUnlock()
		return l
	}
	limitersMu.RUnlock()

	// Buat baru
	limitersMu.Lock()
	defer limitersMu.Unlock()

	// Double check setelah lock
	if l, ok := limiters[name]; ok {
		return l
	}

	l := limiter.New(newStore(), limiter.Rate{
		Period: period,
		Limit:  limit,
	})
	limiters[name] = l
	return l
}

// Handler generik rate limit
func rateLimitHandler(name string, period time.Duration, limit int64, message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := getLimiter(name, period, limit)

		// Key berdasarkan IP
		key := c.ClientIP()

		ctx, err := l.Get(c.Request.Context(), key)
		if err != nil {
			// Jika Redis error, jangan block request — log saja
			fmt.Println("⚠️  Rate limiter error:", err)
			c.Next()
			return
		}

		// Set header informatif
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", ctx.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", ctx.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", ctx.Reset))

		if ctx.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       message,
				"retry_after": ctx.Reset,
			})
			return
		}

		c.Next()
	}
}

// ─── PRESET RATE LIMITERS ────────────────────────────────────

// Global — 300 request per menit per IP
func GlobalRateLimit() gin.HandlerFunc {
	return rateLimitHandler(
		"global",
		1*time.Minute, 300,
		"Terlalu banyak request, coba lagi dalam 1 menit",
	)
}

// Auth — 10 request per menit (login & register)
func AuthRateLimit() gin.HandlerFunc {
	return rateLimitHandler(
		"auth",
		1*time.Minute, 10,
		"Terlalu banyak percobaan login, coba lagi dalam 1 menit",
	)
}

// Message — 20 request per menit (kirim ucapan)
func MessageRateLimit() gin.HandlerFunc {
	return rateLimitHandler(
		"message",
		1*time.Minute, 20,
		"Terlalu banyak ucapan terkirim, coba lagi dalam 1 menit",
	)
}

// API umum — 60 request per menit
func APIRateLimit() gin.HandlerFunc {
	return rateLimitHandler(
		"api",
		1*time.Minute, 60,
		"Terlalu banyak request, coba lagi dalam 1 menit",
	)
}