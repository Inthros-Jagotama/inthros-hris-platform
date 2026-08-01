// Package mailer menyediakan pengiriman email via SMTP menggunakan
// standard library (net/smtp). Cocok untuk development dengan Mailpit
// (SMTP localhost:1025, UI http://localhost:8025) maupun production
// dengan SMTP relay apa pun (termasuk SMTP dari Resend/SES/Mailgun).
package mailer

import (
	"fmt"
	"mime"
	"net/smtp"
	"strings"
	"time"
)

// Config berisi pengaturan koneksi SMTP.
type Config struct {
	// Host SMTP server (mis. "localhost" untuk Mailpit).
	Host string
	// Port SMTP server (Mailpit default 1025).
	Port int
	// Username & Password untuk SMTP auth (kosongkan untuk Mailpit lokal).
	Username string
	Password string
	// From adalah alamat pengirim (mis. "no-reply@hris.local").
	From string
	// FrontendBaseURL adalah base URL frontend tenant (untuk membangun link).
	// Contoh: "http://localhost:5173"
	FrontendBaseURL string
}

// Mailer mengirim email via SMTP.
type Mailer struct {
	cfg Config
}

// New membuat Mailer baru.
func New(cfg Config) *Mailer {
	return &Mailer{cfg: cfg}
}

// Send mengirim email HTML sederhana.
// to: alamat penerima; subject: judul; bodyHTML: isi email (HTML).
func (m *Mailer) Send(to, subject, bodyHTML string) error {
	if m.cfg.Host == "" {
		return fmt.Errorf("mailer: smtp host is not configured")
	}

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	from := m.cfg.From
	if from == "" {
		from = "no-reply@hris.local"
	}

	header := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": `text/html; charset="utf-8"`,
		"Date":         time.Now().Format(time.RFC1123Z),
	}

	var sb strings.Builder
	for k, v := range header {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", k, mime.QEncoding.Encode("utf-8", v)))
	}
	sb.WriteString("\r\n")
	sb.WriteString(bodyHTML)

	// Mailpit lokal tanpa auth → auth nil.
	var auth smtp.Auth
	if m.cfg.Username != "" || m.cfg.Password != "" {
		auth = smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	}

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(sb.String()))
}

// SetupLink membangun URL untuk halaman set-password frontend.
// Path default: /set-password?token=...&company_id=...
// company_id wajib disertakan agar route publik bisa resolve tenant DB.
func (m *Mailer) SetupLink(token, companyID string) string {
	base := strings.TrimRight(m.cfg.FrontendBaseURL, "/")
	if base == "" {
		base = "http://localhost:5173"
	}
	if companyID != "" {
		return fmt.Sprintf("%s/set-password?token=%s&company_id=%s", base, token, companyID)
	}
	return fmt.Sprintf("%s/set-password?token=%s", base, token)
}
