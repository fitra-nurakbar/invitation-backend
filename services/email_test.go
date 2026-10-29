package services

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

func TestPaymentEmailTemplate(t *testing.T) {
	data := PaymentSuccessData{
		UserName:       "Budi Santoso",
		TemplateName:   "Minimalist Botanical",
		Amount:         150000,
		PaymentMethod:  "QRIS",
		PaymentChannel: "GOPAY",
		InvoiceID:      "INV-2026-0001",
	}

	type templateData struct {
		PaymentSuccessData
		AmountFormatted string
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background: #faf9f6; margin: 0; padding: 0; color: #181716; }
    .container { max-width: 540px; margin: 40px auto; background: #ffffff; border: 1px solid #e4dfd5; padding: 40px 32px; box-sizing: border-box; }
    .header { text-align: center; border-bottom: 1px solid #e4dfd5; padding-bottom: 24px; margin-bottom: 28px; }
    .brand { font-family: Georgia, serif; font-size: 20px; font-weight: 700; letter-spacing: 0.15em; text-transform: uppercase; color: #181716; }
    .title { font-size: 18px; font-weight: 600; color: #181716; margin: 0 0 12px; }
    .content p { font-size: 14px; line-height: 1.6; color: #55514b; margin: 0 0 16px; }
    .detail-box { background: #faf9f6; border: 1px solid #e4dfd5; padding: 20px 24px; margin: 24px 0; }
    .detail-row { display: flex; justify-content: space-between; align-items: center; padding: 10px 0; border-bottom: 1px solid #e4dfd5; font-size: 13.5px; }
    .detail-row:last-child { border-bottom: none; padding-top: 14px; margin-top: 4px; }
    .detail-label { color: #6b665f; }
    .detail-value { color: #181716; font-weight: 500; }
    .detail-value-mono { font-family: 'Courier New', Courier, monospace; font-size: 13px; color: #181716; }
    .detail-total { font-family: Georgia, serif; font-size: 18px; font-weight: 700; color: #181716; }
    .badge { display: inline-block; background: #f3f7f4; border: 1px solid #c2d9c8; color: #275235; padding: 3px 10px; font-size: 11px; font-weight: 700; letter-spacing: 0.05em; text-transform: uppercase; }
    .btn-wrapper { text-align: center; margin: 28px 0; }
    .btn { display: inline-block; background: #8a704c; color: #ffffff !important; padding: 12px 28px; text-decoration: none; border: 1px solid #8a704c; font-size: 13.5px; font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; }
    .footer { border-top: 1px solid #e4dfd5; padding-top: 20px; margin-top: 32px; font-size: 12px; color: #8c877e; text-align: center; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <div class="brand">INVITATION</div>
    </div>
    <div class="content">
      <h2 class="title">Pembayaran Berhasil</h2>
      <p>Halo <strong>{{.UserName}}</strong>,</p>
      <p>Pembayaran Anda telah kami terima. Akses template undangan Anda kini telah aktif dan siap digunakan untuk membuat undangan digital.</p>

      <div class="detail-box">
        <div class="detail-row">
          <span class="detail-label">Template</span>
          <span class="detail-value">{{.TemplateName}}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Metode Pembayaran</span>
          <span class="detail-value">{{.PaymentMethod}} - {{.PaymentChannel}}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Nomor Invoice</span>
          <span class="detail-value-mono">{{.InvoiceID}}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Status Transaksi</span>
          <span class="detail-value"><span class="badge">LUNAS</span></span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Total Pembayaran</span>
          <span class="detail-total">Rp {{.AmountFormatted}}</span>
        </div>
      </div>

      <div class="btn-wrapper">
        <a href="http://localhost:3000/dashboard" class="btn">Buka Dashboard</a>
      </div>

      <p style="font-size: 13px; color: #8c877e; margin-top: 20px;">
        Silakan masuk ke dashboard akun Anda untuk mulai mengonfigurasi detail dan membagikan undangan Anda.
      </p>
    </div>
    <div class="footer">
      <p>Pesan ini dikirimkan secara otomatis oleh Invitation App.</p>
    </div>
  </div>
</body>
</html>`

	tmplParsed, err := template.New("email-payment-test").Parse(tmpl)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	var buf bytes.Buffer
	err = tmplParsed.Execute(&buf, templateData{
		PaymentSuccessData: data,
		AmountFormatted:    formatRupiah(data.Amount),
	})
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, "Rp 150.000") {
		t.Errorf("Expected Rp 150.000 in output, got: %s", result)
	}
	if !strings.Contains(result, "INVITATION") {
		t.Errorf("Expected INVITATION brand in output")
	}
}

