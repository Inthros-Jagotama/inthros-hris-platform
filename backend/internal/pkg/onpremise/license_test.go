package onpremise

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func testLicense() License {
	return License{
		CompanyID:      "11111111-2222-3333-4444-555555555555",
		CompanyName:    "PT Contoh",
		ExpiresAt:      time.Now().AddDate(1, 0, 0),
		AllowedModules: []string{"organization", "employee", "payroll"},
		MaxEmployees:   500,
	}
}

func TestGenerateKeyPairAndSignVerify(t *testing.T) {
	privPEM, pubPEM, err := GenerateKeyPair(0)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	licData, err := SignLicense(privPEM, testLicense())
	if err != nil {
		t.Fatalf("SignLicense: %v", err)
	}

	lic, err := VerifyLicense(string(pubPEM), licData)
	if err != nil {
		t.Fatalf("VerifyLicense: %v", err)
	}
	if lic.CompanyName != "PT Contoh" {
		t.Errorf("CompanyName = %q, want PT Contoh", lic.CompanyName)
	}
	if lic.MaxEmployees != 500 {
		t.Errorf("MaxEmployees = %d, want 500", lic.MaxEmployees)
	}
	if !lic.HasModule("payroll") {
		t.Error("HasModule(payroll) = false, want true")
	}
	if lic.HasModule("setting") {
		t.Error("HasModule(setting) = true, want false")
	}
	if !lic.IsValid() {
		t.Error("IsValid() = false, want true")
	}
}

func TestVerify_Expired(t *testing.T) {
	privPEM, pubPEM, err := GenerateKeyPair(0)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	lic := testLicense()
	lic.ExpiresAt = time.Now().Add(-time.Hour)

	licData, err := SignLicense(privPEM, lic)
	if err != nil {
		t.Fatalf("SignLicense: %v", err)
	}

	_, err = VerifyLicense(string(pubPEM), licData)
	if !errors.Is(err, ErrExpired) {
		t.Errorf("VerifyLicense err = %v, want ErrExpired", err)
	}
}

func TestVerify_Tampered(t *testing.T) {
	privPEM, pubPEM, err := GenerateKeyPair(0)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	licData, err := SignLicense(privPEM, testLicense())
	if err != nil {
		t.Fatalf("SignLicense: %v", err)
	}

	// Tamper deterministik: ubah 1 karakter di tengah nilai base64 "payload"
	// (selalu ada), signature tetap asli → payload berbeda → ErrInvalidSig.
	var orig struct {
		Payload   string `json:"payload"`
		Signature string `json:"signature"`
	}
	if err := json.Unmarshal(licData, &orig); err != nil {
		t.Fatalf("unmarshal lic: %v", err)
	}
	mid := len(orig.Payload) / 2
	// Ganti karakter di posisi tengah dengan karakter yang dijamin berbeda
	// (jika aslinya 'B', gunakan 'C') agar selalu mengubah payload.
	newChar := byte('B')
	if orig.Payload[mid] == 'B' {
		newChar = 'C'
	}
	alteredPayload := orig.Payload[:mid] + string(newChar) + orig.Payload[mid+1:]
	altered := []byte(`{"payload":"` + alteredPayload + `","signature":"` + orig.Signature + `"}`)

	_, err = VerifyLicense(string(pubPEM), altered)
	if !errors.Is(err, ErrInvalidSig) {
		t.Errorf("VerifyLicense err = %v, want ErrInvalidSig", err)
	}
}

func TestVerify_WrongKey(t *testing.T) {
	privPEM, _, err := GenerateKeyPair(0)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	_, otherPub, err := GenerateKeyPair(0)
	if err != nil {
		t.Fatalf("GenerateKeyPair(2): %v", err)
	}

	licData, err := SignLicense(privPEM, testLicense())
	if err != nil {
		t.Fatalf("SignLicense: %v", err)
	}

	_, err = VerifyLicense(string(otherPub), licData)
	if !errors.Is(err, ErrInvalidSig) {
		t.Errorf("VerifyLicense err = %v, want ErrInvalidSig", err)
	}
}

func TestVerify_Malformed(t *testing.T) {
	_, pubPEM, err := GenerateKeyPair(0)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	// Data lisensi bukan JSON → ErrMalformed (public key valid, jadi
	// error muncul saat unmarshal file, bukan saat parse key).
	_, err = VerifyLicense(string(pubPEM), []byte("not a json file"))
	if !errors.Is(err, ErrMalformed) {
		t.Errorf("VerifyLicense err = %v, want ErrMalformed", err)
	}
}

func TestWriteAndReadFile(t *testing.T) {
	privPEM, pubPEM, err := GenerateKeyPair(0)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	path := t.TempDir() + "/test.lic"
	if err := WriteLicenseFile(path, privPEM, testLicense()); err != nil {
		t.Fatalf("WriteLicenseFile: %v", err)
	}

	lic, err := ReadLicenseFile(path, string(pubPEM))
	if err != nil {
		t.Fatalf("ReadLicenseFile: %v", err)
	}
	if !lic.HasModule("organization") {
		t.Error("loaded license missing organization module")
	}
}
