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

func ResetDatabase(db *gorm.DB) error {
	// Hapus semua tabel
	err := db.Migrator().DropTable(
		&models.Message{},      // FK -> Invitation
		&models.Invitation{},   // FK -> User, Template
		&models.UserTemplate{}, // FK -> User, Template
		&models.Order{},        // FK -> User
		&models.Template{},
		&models.User{},
	)
	if err != nil {
		return err
	}

	// Buat ulang tabel
	err = db.AutoMigrate(
		&models.User{},
		&models.Template{},
		&models.Order{},
		&models.UserTemplate{},
		&models.Invitation{},
		&models.Message{},
	)

	return err
}

func Run(db *gorm.DB) {
	fmt.Println("🌱 Menjalankan seeder...")

	seedAdminUser(db) // ← admin dulu sebelum user biasa
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
		Email:     "bisniskeluargaofficial@gmail.com",
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
			Slug:              "elegant-rose",
			Name:              "Elegant Rose",
			Price:             150000,
			IsActive:          true,
			OrderDeadlineDays: 7,
			ActiveDaysAfter:   30,
		},
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Slug:              "modern-minimalist",
			Name:              "Modern Minimalist",
			Price:             100000,
			IsActive:          true,
			OrderDeadlineDays: 5,
			ActiveDaysAfter:   14,
		},
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Slug:              "rustic-garden",
			Name:              "Rustic Garden",
			Price:             200000,
			IsActive:          true,
			OrderDeadlineDays: 10,
			ActiveDaysAfter:   60,
		},
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Slug:              "royal-gold",
			Name:              "Royal Gold",
			Price:             250000,
			IsActive:          true,
			OrderDeadlineDays: 14,
			ActiveDaysAfter:   90,
		},
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			Slug:              "simple-free",
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
			Email:     "fitra0612@gmail.com",
			Name:      "Budi Santoso",
			Role:      models.RoleUser,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0001-000000000003"),
			Email:     "jamalgolden203@gmail.com",
			Name:      "Siti Rahayu",
			Role:      models.RoleUser,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0001-000000000004"),
			Email:     "gangsterhitam38@gmail.com",
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
		"bride": map[string]string{
			"name":          "Ani Lestari",
			"place_of_birth": "Jakarta, 15 Februari 1998",
			"parent":        "Putri dari Bapak Hendra & Ibu Wati",
		},
		"groom": map[string]string{
			"name":          "Budi Santoso",
			"place_of_birth": "Bogor, 20 Januari 1996",
			"parent":        "Putra dari Bapak Slamet & Ibu Ningsih",
		},
		"quotes": map[string]string{
			"quote":       "Dan di antara tanda-tanda (kebesaran)-Nya ialah Dia menciptakan pasangan-pasangan untukmu dari jenismu sendiri, agar kamu cenderung dan merasa tenteram kepadanya, dan Dia menjadikan di antaramu rasa kasih dan sayang.",
			"attribution": "(QS. Ar-Rum: 21)",
		},
		"wedding_event": []map[string]string{
			{
				"name":     "akad",
				"date":     "Kamis, 25 Desember 2025",
				"time":     "08:00 - 10:00 WIB",
				"place":    "Masjid Al-Ikhlas",
				"address":  "Jl. Merdeka No. 1, Jakarta Pusat",
				"maps_url": "https://maps.google.com/?q=-6.2,106.8",
			},
			{
				"name":     "resepsi",
				"date":     "Kamis, 25 Desember 2025",
				"time":     "11:00 - 14:00 WIB",
				"place":    "Gedung Serbaguna Merdeka",
				"address":  "Jl. Merdeka No. 1, Jakarta Pusat",
				"maps_url": "https://maps.google.com/?q=-6.2,106.8",
			},
		},
		"streaming_platform": "https://instagram.com/budiani_wedding",
		"gallery":            map[string]interface{}{},
		"love_story":         map[string]interface{}{},
		"wedding_gift": []map[string]string{
			{
				"platform": "bca",
				"name":     "Budi Santoso",
				"id":       "1234567890",
			},
			{
				"platform": "gopay",
				"name":     "Budi Santoso",
				"id":       "08123456789",
			},
		},
	})

	detailSiti, _ := json.Marshal(map[string]interface{}{
		"bride": map[string]string{
			"name":          "Siti Rahayu",
			"place_of_birth": "Bogor, 10 Maret 1999",
			"parent":        "Putri dari Bapak Suparman & Ibu Surati",
		},
		"groom": map[string]string{
			"name":          "Rizky Pratama",
			"place_of_birth": "Bandung, 5 Juli 1997",
			"parent":        "Putra dari Bapak Darmawan & Ibu Yuli",
		},
		"quotes": map[string]string{
			"quote":       "Maka nikmat Tuhanmu yang manakah yang kamu dustakan?",
			"attribution": "(QS. Ar-Rahman: 13)",
		},
		"wedding_event": []map[string]string{
			{
				"name":     "akad",
				"date":     "Sabtu, 14 Februari 2026",
				"time":     "09:00 - 11:00 WIB",
				"place":    "Rumah Mempelai Wanita",
				"address":  "Jl. Bukit Indah No. 5, Bogor",
				"maps_url": "https://maps.google.com/?q=-6.5,106.8",
			},
			{
				"name":     "resepsi",
				"date":     "Sabtu, 14 Februari 2026",
				"time":     "12:00 - 15:00 WIB",
				"place":    "Villa Bukit Indah",
				"address":  "Jl. Bukit Indah No. 10, Bogor",
				"maps_url": "https://maps.google.com/?q=-6.5,106.8",
			},
		},
		"streaming_platform": "https://youtube.com/@sitirizkywedding",
		"gallery":            map[string]interface{}{},
		"love_story":         map[string]interface{}{},
		"wedding_gift": []map[string]string{
			{
				"platform": "bni",
				"name":     "Siti Rahayu",
				"id":       "0987654321",
			},
		},
	})

	detailAhmad, _ := json.Marshal(map[string]interface{}{
		"bride": map[string]string{
			"name":          "Dewi Anggraini",
			"place_of_birth": "Bandung, 3 Agustus 2000",
			"parent":        "Putri dari Bapak Wahyu & Ibu Endah",
		},
		"groom": map[string]string{
			"name":          "Ahmad Fauzi",
			"place_of_birth": "Surabaya, 17 November 1997",
			"parent":        "Putra dari Bapak Ridwan & Ibu Fatimah",
		},
		"quotes": map[string]string{
			"quote":       "Wahai manusia! Sungguh, Kami telah menciptakan kamu dari seorang laki-laki dan seorang perempuan, kemudian Kami jadikan kamu berbangsa-bangsa dan bersuku-suku agar kamu saling mengenal.",
			"attribution": "(QS. Al-Hujurat: 13)",
		},
		"wedding_event": []map[string]string{
			{
				"name":     "akad",
				"date":     "Minggu, 8 Maret 2026",
				"time":     "07:30 - 09:00 WIB",
				"place":    "Masjid Agung Bandung",
				"address":  "Jl. Asia Afrika, Bandung",
				"maps_url": "https://maps.google.com/?q=-6.9,107.6",
			},
			{
				"name":     "resepsi",
				"date":     "Minggu, 8 Maret 2026",
				"time":     "10:00 - 13:00 WIB",
				"place":    "Hotel Bintang Lima Bandung",
				"address":  "Jl. Braga No. 8, Bandung",
				"maps_url": "https://maps.google.com/?q=-6.9,107.6",
			},
		},
		"streaming_platform": "",
		"gallery":            map[string]interface{}{},
		"love_story":         map[string]interface{}{},
		"wedding_gift": []map[string]string{
			{
				"platform": "mandiri",
				"name":     "Ahmad Fauzi",
				"id":       "1122334455",
			},
			{
				"platform": "bri",
				"name":     "Dewi Anggraini",
				"id":       "5544332211",
			},
		},
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
		// Cek apakah sudah ada
		var existing models.Invitation
		if db.Where("id = ?", inv.ID).First(&existing).Error == nil {
			// Update detail dengan struktur baru
			db.Model(&existing).Update("detail", inv.Detail)
			fmt.Println("  🔄 Updated invitation:", inv.Slug)
			continue
		}

		result := db.Create(&inv)
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
