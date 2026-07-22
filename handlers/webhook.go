package handlers

import (
	"fmt"
	"invitation-app/config"
	"invitation-app/models"
	"invitation-app/services"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type XenditWebhookPayload struct {
	ID             string  `json:"id"`
	ExternalID     string  `json:"external_id"`
	Status         string  `json:"status"`
	Amount         float64 `json:"amount"`
	PaymentMethod  string  `json:"payment_method"`
	PaymentChannel string  `json:"payment_channel"`
	PaidAt         string  `json:"paid_at"`
}

// POST /webhook/xendit
func XenditWebhook(c *gin.Context) {
	// Verifikasi token webhook dari header
	token := c.GetHeader("x-callback-token")
	if token != os.Getenv("XENDIT_WEBHOOK_TOKEN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	var payload XenditWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"})
		return
	}

	// Cari order berdasarkan invoice_id
	var order models.Order
	result := config.DB.
		Preload("User").
		Preload("Template").
		Where("invoice_id = ?", payload.ID).
		First(&order)

	if result.Error != nil {
		// Return 200 supaya Xendit tidak retry terus
		c.JSON(http.StatusOK, gin.H{"message": "Order tidak ditemukan, diabaikan"})
		return
	}

	switch payload.Status {

	case "PAID":
		if order.Status == models.OrderStatusPaid {
			c.JSON(http.StatusOK, gin.H{"message": "Sudah diproses"})
			return
		}

		paidAt := time.Now()

		// 1. Update status order
		if err := config.DB.Model(&order).Updates(map[string]interface{}{
			"status":          models.OrderStatusPaid,
			"paid_at":         paidAt,
			"payment_method":  payload.PaymentMethod,
			"payment_channel": payload.PaymentChannel,
			"updated_at":      time.Now(),
		}).Error; err != nil {
			fmt.Println("❌ Gagal update order:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update order"})
			return
		}

		// 2. Aktifkan akses template untuk user
		if err := templateAccessSvc().Grant(
			order.UserID,
			order.TemplateID,
			order.ID,
		); err != nil {
			fmt.Println("❌ Gagal grant akses template:", err)
			// Tidak return — tetap kirim email
		} else {
			fmt.Printf("✅ Akses template %s diberikan ke user %s\n",
				order.Template.Name, order.User.Name)
		}

		// 3. Kirim email notifikasi ke user
		go func() {
			err := emailSvc().SendPaymentSuccess(order.User.Email, services.PaymentSuccessData{
				UserName:       order.User.Name,
				TemplateName:   order.Template.Name,
				Amount:         order.Amount,
				PaymentMethod:  payload.PaymentMethod,
				PaymentChannel: payload.PaymentChannel,
				InvoiceID:      payload.ID,
			})
			if err != nil {
				fmt.Println("❌ Gagal kirim email:", err)
			} else {
				fmt.Printf("✅ Email notifikasi terkirim ke %s\n", order.User.Email)
			}
		}()

	case "EXPIRED":
		config.DB.Model(&order).Updates(map[string]interface{}{
			"status":     models.OrderStatusExpired,
			"updated_at": time.Now(),
		})
		fmt.Printf("⏰ Order %s expired\n", order.ID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook diproses"})
}
