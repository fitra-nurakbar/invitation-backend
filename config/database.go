package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL tidak ditemukan di .env")
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Tambah prefer_simple_protocol untuk fix error prepared statement di Supabase
	dsn = dsn + "?prefer_simple_protocol=true"

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // ← fix error "prepared statement already exists"
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// Connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Gagal get sql.DB:", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("✅ Koneksi ke Supabase berhasil!")
	DB = db
}