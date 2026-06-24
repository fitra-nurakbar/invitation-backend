package main

import (
	"fmt"
	"invitation-app/config"
	"invitation-app/routes"
	"invitation-app/seeders"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️  .env tidak ditemukan, menggunakan environment variables sistem")
	}

	requiredEnvs := []string{
		"DATABASE_URL",
		"JWT_SECRET",
		"XENDIT_SECRET_KEY",
		"XENDIT_WEBHOOK_TOKEN",
	}
	for _, env := range requiredEnvs {
		val := os.Getenv(env)
		if val == "" {
			log.Fatalf("❌ Environment variable %s tidak ditemukan!", env)
		}
		// Tampilkan 12 karakter pertama saja
		preview := val
		if len(val) > 12 {
			preview = val[:12]
		}
		fmt.Printf("✅ %s = %s...\n", env, preview)
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	config.InitGoogleOAuth()
	config.ConnectRedis()
	config.RunMigrations()
	config.ConnectDatabase()

	if len(os.Args) > 1 && os.Args[1] == "--seed" {
		seeders.Run(config.DB)
		fmt.Println("🎉 Seeding selesai.")
		return
	}

	r := gin.Default()
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
	r.Run(":" + port)
}
