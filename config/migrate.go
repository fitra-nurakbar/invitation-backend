package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() {
	databaseURL := os.Getenv("DATABASE_URL")

	migrationsURL := getMigrationsURL()
	fmt.Println("📁 Migrations path:", migrationsURL)

	m, err := migrate.New(migrationsURL, databaseURL)
	if err != nil {
		log.Fatal("❌ Gagal inisialisasi migrate:", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("✅ Migration tidak ada perubahan.")
			return
		}
		log.Fatal("❌ Gagal menjalankan migration:", err)
	}

	fmt.Println("✅ Migration selesai!")
}

func getMigrationsURL() string {
	// Cek apakah jalan di Docker (path absolut /app/migrations)
	dockerPath := "/app/migrations"
	if _, err := os.Stat(dockerPath); err == nil {
		return fmt.Sprintf("file://%s", dockerPath)
	}

	// Local development — relatif ke file ini
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(filename), "..")
	migrationsPath, _ := filepath.Abs(filepath.Join(rootDir, "migrations"))
	migrationsSlash := filepath.ToSlash(migrationsPath)

	return fmt.Sprintf("file://%s", migrationsSlash)
}