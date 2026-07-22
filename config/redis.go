package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("❌ Gagal parse REDIS_URL:", err)
	}

	// Override password & DB jika ada di env
	if pass := os.Getenv("REDIS_PASSWORD"); pass != "" {
		opt.Password = pass
	}

	client := redis.NewClient(opt)

	// Test koneksi
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("❌ Gagal koneksi ke Redis:", err)
	}

	fmt.Println("✅ Koneksi ke Redis berhasil!")
	RedisClient = client
}
