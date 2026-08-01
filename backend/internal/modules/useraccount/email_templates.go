package useraccount

import "fmt"

const subjectSetupAccount = "HRIS — Set Password Akun Anda"

// bodySetupAccount membuat body email HTML berisi link set-password.
func bodySetupAccount(employeeName, link string) string {
	greeting := "Halo,"
	if employeeName != "" {
		greeting = fmt.Sprintf("Halo, %s,", employeeName)
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family:Arial,Helvetica,sans-serif;background:#f4f6f8;margin:0;padding:24px;">
  <div style="max-width:520px;margin:0 auto;background:#ffffff;border-radius:8px;overflow:hidden;border:1px solid #e5e7eb;">
    <div style="background:#059669;padding:20px 24px;">
      <span style="color:#ffffff;font-size:18px;font-weight:bold;">HRIS</span>
    </div>
    <div style="padding:24px;">
      <p style="font-size:15px;color:#111827;">%s</p>
      <p style="font-size:14px;color:#374151;">Akun Anda telah dibuat. Silakan klik tombol di bawah untuk membuat password dan mulai menggunakan aplikasi.</p>
      <p style="text-align:center;margin:28px 0;">
        <a href="%s" style="background:#059669;color:#ffffff;text-decoration:none;padding:12px 28px;border-radius:6px;font-size:14px;font-weight:bold;display:inline-block;">Set Password</a>
      </p>
      <p style="font-size:13px;color:#6b7280;">Link ini berlaku selama 48 jam dan hanya bisa digunakan sekali.</p>
      <p style="font-size:13px;color:#6b7280;margin-top:8px;">Jika tombol tidak berfungsi, salin tautan berikut ke browser Anda:<br>
        <span style="color:#059669;word-break:break-all;">%s</span>
      </p>
    </div>
  </div>
</body>
</html>`, greeting, link, link)
}
