package handlers

import (
	"fmt"
	"invitation-app/config"
	"invitation-app/models"
	"invitation-app/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// POST /orders — buat order baru
func CreateOrder(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var input struct {
		TemplateID string `json:"template_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templateID, err := uuid.Parse(input.TemplateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template_id tidak valid"})
		return
	}

	// Ambil data user
	var user models.User
	if result := config.DB.First(&user, "id = ?", userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Ambil data template
	var template models.Template
	if result := config.DB.First(&template, "id = ? AND is_active = true", templateID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template tidak ditemukan atau tidak aktif"})
		return
	}

	// Cek apakah user sudah punya order pending untuk template ini
	var existingOrder models.Order
	result := config.DB.Where(
		"user_id = ? AND template_id = ? AND status = ?",
		userID, templateID, models.OrderStatusPending,
	).First(&existingOrder)

	if result.Error == nil {
		// Order pending sudah ada — kembalikan order yang sama
		c.JSON(http.StatusOK, gin.H{
			"message": "Order sudah ada, silakan selesaikan pembayaran",
			"data":    existingOrder,
		})
		return
	}

	// Cek apakah user sudah pernah beli template ini
	var paidOrder models.Order
	result = config.DB.Where(
		"user_id = ? AND template_id = ? AND status = ?",
		userID, templateID, models.OrderStatusPaid,
	).First(&paidOrder)

	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Kamu sudah memiliki template ini"})
		return
	}

	// Buat order ID dulu
	orderID := uuid.New()

	// Buat invoice di Xendit
invoiceResult, err := xenditSvc().CreateInvoice(services.CreateInvoiceParams{
    OrderID:  orderID,
    User:     user,
    Template: template,
})
if err != nil {
    fmt.Println("❌ Xendit error:", err) // ← tambah ini
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat invoice pembayaran", "detail": err.Error()})
    return
}

	// Simpan order ke database
	order := models.Order{
		ID:         orderID,
		UserID:     userID,
		TemplateID: templateID,
		InvoiceID:  invoiceResult.InvoiceID,
		InvoiceURL: invoiceResult.InvoiceURL,
		Amount:     invoiceResult.Amount,
		Status:     models.OrderStatusPending,
		ExpiresAt:  invoiceResult.ExpiresAt,
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Order berhasil dibuat, silakan selesaikan pembayaran",
		"data":        order,
		"invoice_url": invoiceResult.InvoiceURL,
	})
}

// GET /orders — list order milik user
func GetOrders(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var orders []models.Order
	config.DB.
		Preload("Template").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders)

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// GET /orders/:id — detail order
func GetOrder(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var order models.Order
	result := config.DB.
		Preload("Template").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// POST /orders/:id/cancel — batalkan order pending
func CancelOrder(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var order models.Order
	result := config.DB.Where(
		"id = ? AND user_id = ? AND status = ?",
		orderID, userID, models.OrderStatusPending,
	).First(&order)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan atau tidak bisa dibatalkan"})
		return
	}

	// Expire invoice di Xendit
	if err := xenditSvc().ExpireInvoice(order.InvoiceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membatalkan invoice"})
		return
	}

	// Update status order
	now := time.Now()
	config.DB.Model(&order).Updates(map[string]interface{}{
		"status":     models.OrderStatusExpired,
		"updated_at": now,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Order berhasil dibatalkan"})
}

// GET /my-templates — template yang sudah dibeli user
func GetMyTemplates(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	accesses, err := templateAccessSvc().GetUserTemplates(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": accesses})
}

// ─── ADMIN ───────────────────────────────────────────────────

// GET /admin/orders — semua order (admin)
func GetAllOrders(c *gin.Context) {
	var orders []models.Order
	config.DB.
		Preload("User").
		Preload("Template").
		Order("created_at DESC").
		Find(&orders)

	c.JSON(http.StatusOK, gin.H{"data": orders})
}
