package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	host     string
	port     int
	email    string
	password string
	sender   string
}

type PaymentSuccessData struct {
	UserName       string
	TemplateName   string
	Amount         int
	PaymentMethod  string
	PaymentChannel string
	InvoiceID      string
}

func NewEmailService() *EmailService {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return &EmailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     port,
		email:    os.Getenv("SMTP_EMAIL"),
		password: os.Getenv("SMTP_PASSWORD"),
		sender:   os.Getenv("SMTP_SENDER_NAME"),
	}
}

// send — kirim email generik
func (s *EmailService) send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.sender, s.email))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.host, s.port, s.email, s.password)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("gagal kirim email: %w", err)
	}
	return nil
}

// SendPaymentSuccess — notifikasi pembayaran berhasil
func (s *EmailService) SendPaymentSuccess(to string, data PaymentSuccessData) error {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 40px auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
    .header { background: #4F46E5; padding: 32px; text-align: center; }
    .header h1 { color: #fff; margin: 0; font-size: 24px; }
    .body { padding: 32px; }
    .body p { color: #555; line-height: 1.6; }
    .detail-box { background: #F9FAFB; border-radius: 8px; padding: 20px; margin: 24px 0; }
    .detail-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #E5E7EB; }
    .detail-row:last-child { border-bottom: none; font-weight: bold; }
    .detail-label { color: #6B7280; }
    .detail-value { color: #111827; }
    .badge { display: inline-block; background: #D1FAE5; color: #065F46; padding: 4px 12px; border-radius: 999px; font-size: 13px; font-weight: bold; }
    .footer { background: #F9FAFB; padding: 20px; text-align: center; }
    .footer p { color: #9CA3AF; font-size: 12px; margin: 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>✅ Pembayaran Berhasil!</h1>
    </div>
    <div class="body">
      <p>Halo <strong>{{.UserName}}</strong>,</p>
      <p>Pembayaran kamu telah kami terima. Template undangan kamu sudah aktif dan siap digunakan!</p>

      <div class="detail-box">
        <div class="detail-row">
          <span class="detail-label">Template</span>
          <span class="detail-value">{{.TemplateName}}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Metode Bayar</span>
          <span class="detail-value">{{.PaymentMethod}} - {{.PaymentChannel}}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Invoice ID</span>
          <span class="detail-value">{{.InvoiceID}}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Status</span>
          <span class="detail-value"><span class="badge">LUNAS</span></span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Total Bayar</span>
          <span class="detail-value">Rp {{.AmountFormatted}}</span>
        </div>
      </div>

      <p>Silakan login ke dashboard untuk mulai membuat undangan digitalmu. 🎉</p>
    </div>
    <div class="footer">
      <p>Email ini dikirim otomatis, mohon tidak membalas.</p>
    </div>
  </div>
</body>
</html>`

	// Tambah field format rupiah
	type templateData struct {
		PaymentSuccessData
		AmountFormatted string
	}

	formatted := formatRupiah(data.Amount)
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("gagal parse template email: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, templateData{
		PaymentSuccessData: data,
		AmountFormatted:    formatted,
	}); err != nil {
		return fmt.Errorf("gagal render template email: %w", err)
	}

	return s.send(to, "✅ Pembayaran Berhasil - "+data.TemplateName, buf.String())
}

// formatRupiah — format angka ke format Rupiah
func formatRupiah(amount int) string {
	s := strconv.Itoa(amount)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	return result
}

// SendMagicLink — kirim email magic link
func (s *EmailService) SendMagicLink(to, name, magicURL string) error {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 40px auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
    .header { background: #4F46E5; padding: 32px; text-align: center; }
    .header h1 { color: #fff; margin: 0; font-size: 24px; }
    .body { padding: 32px; }
    .body p { color: #555; line-height: 1.6; }
    .btn { display: block; width: fit-content; margin: 24px auto; background: #4F46E5; color: #fff !important; padding: 14px 32px; border-radius: 8px; text-decoration: none; font-weight: bold; font-size: 16px; }
    .note { background: #FEF3C7; border-radius: 8px; padding: 16px; margin-top: 24px; }
    .note p { color: #92400E; margin: 0; font-size: 14px; }
    .footer { background: #F9FAFB; padding: 20px; text-align: center; }
    .footer p { color: #9CA3AF; font-size: 12px; margin: 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🔐 Login ke Invitation App</h1>
    </div>
    <div class="body">
      <p>Halo <strong>{{.Name}}</strong>,</p>
      <p>Kamu meminta link untuk masuk ke akun Invitation App. Klik tombol di bawah untuk login:</p>

      <a href="{{.MagicURL}}" class="btn">Login Sekarang</a>

      <div class="note">
        <p>⏰ Link ini hanya berlaku selama <strong>15 menit</strong> dan hanya bisa digunakan <strong>sekali</strong>.</p>
      </div>

      <p style="margin-top:24px; color:#9CA3AF; font-size:13px;">
        Jika kamu tidak meminta link ini, abaikan email ini. Akun kamu tetap aman.
      </p>
    </div>
    <div class="footer">
      <p>Email ini dikirim otomatis, mohon tidak membalas.</p>
    </div>
  </div>
</body>
</html>`

	type templateData struct {
		Name     string
		MagicURL string
	}

	t, err := template.New("magic-link").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("gagal parse template magic link: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, templateData{
		Name:     name,
		MagicURL: magicURL,
	}); err != nil {
		return fmt.Errorf("gagal render template magic link: %w", err)
	}

	return s.send(to, "🔐 Link Login Invitation App", buf.String())
}