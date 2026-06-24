package services

import (
	"context"
	"fmt"
	"invitation-app/models"
	"os"
	"time"

	"github.com/google/uuid"
	xendit "github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/invoice"
)

type XenditService struct {
	client *xendit.APIClient
}

type CreateInvoiceParams struct {
	OrderID  uuid.UUID
	User     models.User
	Template models.Template
}

type InvoiceResult struct {
	InvoiceID  string
	InvoiceURL string
	ExpiresAt  time.Time
	Amount     int
}

func NewXenditService() *XenditService {
	secretKey := os.Getenv("XENDIT_SECRET_KEY")

	// Validasi key tidak kosong
	if secretKey == "" {
		fmt.Println("⚠️  XENDIT_SECRET_KEY tidak ditemukan di environment!")
	} else {
		fmt.Printf("✅ Xendit key loaded: %s...\n", secretKey[:12]) // tampilkan 12 karakter pertama saja
	}

	client := xendit.NewClient(secretKey)
	return &XenditService{client: client}
}

func strPtr(s string) *string {
	return &s
}

func (s *XenditService) CreateInvoice(params CreateInvoiceParams) (*InvoiceResult, error) {
	secretKey := os.Getenv("XENDIT_SECRET_KEY")

	// Double check saat runtime — karena NewXenditService() dipanggil saat app start
	// tapi .env mungkin belum ter-load saat itu
	if secretKey == "" {
		return nil, fmt.Errorf("XENDIT_SECRET_KEY kosong")
	}

	// Re-init client dengan key terbaru untuk memastikan key ter-load
	s.client = xendit.NewClient(secretKey)

	externalID := fmt.Sprintf("order-%s", params.OrderID.String())
	expiresAt  := time.Now().Add(24 * time.Hour)
	amount     := float64(params.Template.Price)
	desc       := fmt.Sprintf("Pembelian Template %s - %s", params.Template.Name, params.User.Name)
	duration   := "86400"
	currency   := "IDR"
	successURL := os.Getenv("XENDIT_SUCCESS_URL")
	failureURL := os.Getenv("XENDIT_FAILURE_URL")

	req := invoice.NewCreateInvoiceRequest(externalID, amount)
	req.Description        = strPtr(desc)
	req.PayerEmail         = strPtr(params.User.Email)
	req.Currency           = strPtr(currency)
	req.InvoiceDuration    = strPtr(duration)
	req.SuccessRedirectUrl = strPtr(successURL)
	req.FailureRedirectUrl = strPtr(failureURL)
	req.PaymentMethods     = []string{
		"BCA", "BNI", "BRI", "MANDIRI", "PERMATA",
		"QRIS",
		"OVO", "DANA", "LINKAJA", "SHOPEEPAY",
		"ALFAMART", "INDOMARET",
	}

	resp, _, err := s.client.InvoiceApi.
		CreateInvoice(context.Background()).
		CreateInvoiceRequest(*req).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("gagal membuat invoice Xendit: %w", err)
	}

	return &InvoiceResult{
		InvoiceID:  resp.GetId(),
		InvoiceURL: resp.GetInvoiceUrl(),
		ExpiresAt:  expiresAt,
		Amount:     params.Template.Price,
	}, nil
}

func (s *XenditService) GetInvoice(invoiceID string) (*invoice.Invoice, error) {
	resp, _, err := s.client.InvoiceApi.
		GetInvoiceById(context.Background(), invoiceID).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("gagal get invoice: %w", err)
	}
	return resp, nil
}

func (s *XenditService) ExpireInvoice(invoiceID string) error {
	_, _, err := s.client.InvoiceApi.
		ExpireInvoice(context.Background(), invoiceID).
		Execute()
	if err != nil {
		return fmt.Errorf("gagal expire invoice: %w", err)
	}
	return nil
}