package seeders

import (
	"encoding/json"
	"fmt"
	"invitation-app/models"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	fmt.Println("🌱 Menjalankan seeder...")

	seedAdminUser(db)   // ← admin dulu sebelum user biasa
	seedTemplates(db)
	seedUsers(db)
	seedInvitations(db)
	seedMessages(db)

	fmt.Println("✅ Seeder selesai!")
}

// ─── ADMIN USER ──────────────────────────────────────────────
func seedAdminUser(db *gorm.DB) {
	adminID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	var existing models.User
	if db.Where("id = ?", adminID).First(&existing).Error == nil {
		fmt.Println("  ⏭ Admin sudah ada, skip.")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin@12345"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Gagal hash password admin:", err)
	}

	pwd := string(hashedPassword) // ← buat variable dulu
	admin := models.User{
		ID:        adminID,
		Email:     "admin@invitation.com",
		Name:      "Super Admin",
		Password:  &pwd, // ← pointer
		Role:      models.RoleAdmin,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Println("❌ Error seed admin:", err)
		return
	}

	fmt.Println("  ✔ Admin user dibuat")
	fmt.Println("     Email    :", admin.Email)
	fmt.Println("     Password : Admin@12345")
}

// ─── TEMPLATES ───────────────────────────────────────────────
func seedTemplates(db *gorm.DB) {
	templates := []models.Template{
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Name:              "Elegant Rose",
			Price:             150000,
			IsActive:          true,
			OrderDeadlineDays: 7,
			ActiveDaysAfter:   30,
		},
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Name:              "Modern Minimalist",
			Price:             100000,
			IsActive:          true,
			OrderDeadlineDays: 5,
			ActiveDaysAfter:   14,
		},
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Name:              "Rustic Garden",
			Price:             200000,
			IsActive:          true,
			OrderDeadlineDays: 10,
			ActiveDaysAfter:   60,
		},
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Name:              "Royal Gold",
			Price:             250000,
			IsActive:          true,
			OrderDeadlineDays: 14,
			ActiveDaysAfter:   90,
		},
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			Name:              "Simple Free",
			Price:             0,
			IsActive:          true,
			OrderDeadlineDays: 3,
			ActiveDaysAfter:   7,
		},
	}

	for _, t := range templates {
		result := db.Where("id = ?", t.ID).FirstOrCreate(&t)
		if result.Error != nil {
			log.Println("❌ Error seed template:", result.Error)
		} else {
			fmt.Println("  ✔ Template:", t.Name)
		}
	}
}

// ─── USERS ───────────────────────────────────────────────────
func seedUsers(db *gorm.DB) {
	users := []models.User{
		{
			ID:        uuid.MustParse("00000000-0000-0000-0001-000000000002"),
			Email:     "budi.santoso@gmail.com",
			Name:      "Budi Santoso",
			Role:      models.RoleUser,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0001-000000000003"),
			Email:     "siti.rahayu@gmail.com",
			Name:      "Siti Rahayu",
			Role:      models.RoleUser,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0001-000000000004"),
			Email:     "ahmad.fauzi@gmail.com",
			Name:      "Ahmad Fauzi",
			Role:      models.RoleUser,
			CreatedAt: time.Now(),
		},
	}

	for _, u := range users {
		result := db.Where("id = ?", u.ID).FirstOrCreate(&u)
		if result.Error != nil {
			log.Println("❌ Error seed user:", result.Error)
		} else {
			fmt.Println("  ✔ User:", u.Name)
		}
	}

	fmt.Println("     Password semua user : User@12345")
}

// ─── INVITATIONS ─────────────────────────────────────────────
func seedInvitations(db *gorm.DB) {
	detailBudi, _ := json.Marshal(map[string]interface{}{
		"bride":          "Ani Lestari",
		"groom":          "Budi Santoso",
		"venue":          "Gedung Serbaguna Merdeka, Jakarta",
		"venue_maps_url": "https://maps.google.com/?q=-6.2,106.8",
		"akad": map[string]string{
			"time":  "08:00 WIB",
			"place": "Masjid Al-Ikhlas Jakarta",
		},
		"resepsi": map[string]string{
			"time":  "11:00 WIB",
			"place": "Gedung Serbaguna Merdeka",
		},
		"bride_parents": "Bapak Hendra & Ibu Wati",
		"groom_parents": "Bapak Slamet & Ibu Ningsih",
		"music":         "Sempurna - Andra and The Backbone",
	})

	detailSiti, _ := json.Marshal(map[string]interface{}{
		"bride":          "Siti Rahayu",
		"groom":          "Rizky Pratama",
		"venue":          "Villa Bukit Indah, Bogor",
		"venue_maps_url": "https://maps.google.com/?q=-6.5,106.8",
		"akad": map[string]string{
			"time":  "09:00 WIB",
			"place": "Rumah Mempelai Wanita",
		},
		"resepsi": map[string]string{
			"time":  "12:00 WIB",
			"place": "Villa Bukit Indah",
		},
		"bride_parents": "Bapak Suparman & Ibu Surati",
		"groom_parents": "Bapak Darmawan & Ibu Yuli",
		"music":         "Untuk Selamanya - Sheila on 7",
	})

	detailAhmad, _ := json.Marshal(map[string]interface{}{
		"bride":          "Dewi Anggraini",
		"groom":          "Ahmad Fauzi",
		"venue":          "Hotel Bintang Lima, Bandung",
		"venue_maps_url": "https://maps.google.com/?q=-6.9,107.6",
		"akad": map[string]string{
			"time":  "07:30 WIB",
			"place": "Masjid Agung Bandung",
		},
		"resepsi": map[string]string{
			"time":  "10:00 WIB",
			"place": "Hotel Bintang Lima Bandung",
		},
		"bride_parents": "Bapak Wahyu & Ibu Endah",
		"groom_parents": "Bapak Ridwan & Ibu Fatimah",
		"music":         "Cinta Sejati - Bunga Citra Lestari",
	})

	invitations := []models.Invitation{
		{
			ID:         uuid.MustParse("00000000-0000-0000-0002-000000000001"),
			UserID:     uuid.MustParse("00000000-0000-0000-0001-000000000002"),
			TemplateID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Slug:       "pernikahan-budi-ani",
			EventDate:  time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
			Status:     models.StatusActive,
			ExpiresAt:  time.Date(2026, 1, 24, 0, 0, 0, 0, time.UTC),
			Detail:     datatypes.JSON(detailBudi),
		},
		{
			ID:         uuid.MustParse("00000000-0000-0000-0002-000000000002"),
			UserID:     uuid.MustParse("00000000-0000-0000-0001-000000000003"),
			TemplateID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Slug:       "pernikahan-siti-rizky",
			EventDate:  time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
			Status:     models.StatusActive,
			ExpiresAt:  time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
			Detail:     datatypes.JSON(detailSiti),
		},
		{
			ID:         uuid.MustParse("00000000-0000-0000-0002-000000000003"),
			UserID:     uuid.MustParse("00000000-0000-0000-0001-000000000004"),
			TemplateID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Slug:       "pernikahan-ahmad-dewi",
			EventDate:  time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC),
			Status:     models.StatusDraft,
			ExpiresAt:  time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC),
			Detail:     datatypes.JSON(detailAhmad),
		},
	}

	for _, inv := range invitations {
		result := db.Where("id = ?", inv.ID).FirstOrCreate(&inv)
		if result.Error != nil {
			log.Println("❌ Error seed invitation:", result.Error)
		} else {
			fmt.Println("  ✔ Invitation:", inv.Slug)
		}
	}
}

// ─── MESSAGES ────────────────────────────────────────────────
func seedMessages(db *gorm.DB) {
	messages := []models.Message{
		{
			ID:           uuid.New(),
			InvitationID: uuid.MustParse("00000000-0000-0000-0002-000000000001"),
			Name:         "Pak Hendra",
			Message:      "Selamat menempuh hidup baru! Semoga menjadi keluarga yang sakinah, mawaddah, warahmah.",
			IPAddress:    "127.0.0.1",
		},
		{
			ID:           uuid.New(),
			InvitationID: uuid.MustParse("00000000-0000-0000-0002-000000000001"),
			Name:         "Bu Sari",
			Message:      "Barakallah! Semoga pernikahan kalian langgeng hingga akhir hayat.",
			IPAddress:    "192.168.1.1",
		},
		{
			ID:           uuid.New(),
			InvitationID: uuid.MustParse("00000000-0000-0000-0002-000000000001"),
			Name:         "Teman Kantor",
			Message:      "Congrats Bud! Akhirnya nikah juga haha. Semoga bahagia selalu ya!",
			IPAddress:    "10.0.0.1",
		},
		{
			ID:           uuid.New(),
			InvitationID: uuid.MustParse("00000000-0000-0000-0002-000000000002"),
			Name:         "Keluarga Besar",
			Message:      "Selamat ya Siti! Semoga menjadi keluarga yang harmonis dan bahagia.",
			IPAddress:    "127.0.0.1",
		},
		{
			ID:           uuid.New(),
			InvitationID: uuid.MustParse("00000000-0000-0000-0002-000000000002"),
			Name:         "Sahabat SMA",
			Message:      "Waaah akhirnya! Selamat pengantin baru, semoga lekas punya momongan ya!",
			IPAddress:    "192.168.1.2",
		},
	}

	for _, m := range messages {
		result := db.Create(&m)
		if result.Error != nil {
			log.Println("❌ Error seed message:", result.Error)
		} else {
			fmt.Println("  ✔ Message dari:", m.Name)
		}
	}
}