// Package onpremise menyediakan mesin lisensi On-Premise berbasis RSA.
//
// Format file .lic (JSON, dua field):
//
//	{
//	  "payload":   "<base64 dari JSON License>",
//	  "signature": "<base64 dari RSA-SHA256 signature atas payload>"
//	}
//
// Alur:
//  1. Platform Admin meng-generate RSA keypair (private + public PEM) via CLI `licensectl`.
//  2. Platform Admin meng-generate file `.lic` dengan private key (berisi
//     expires_at, allowed_modules, max_employees).
//  3. Deployment On-Premise membaca file `.lic` + public key saat startup
//     (DEPLOYMENT_MODE=ON_PREMISE), memverifikasi signature & expiry, lalu
//     daftar modul yang diizinkan dipakai sebagai sumber lisensi modul
//     (pengganti company_modules DB pada mode SaaS).
package onpremise

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"
)

// License adalah payload lisensi On-Premise.
type License struct {
	CompanyID      string    `json:"company_id"`
	CompanyName    string    `json:"company_name"`
	ExpiresAt      time.Time `json:"expires_at"`
	AllowedModules []string  `json:"allowed_modules"`
	MaxEmployees   int       `json:"max_employees"`
}

// licFile adalah struktur file .lic di disk.
type licFile struct {
	Payload   string `json:"payload"`   // base64 dari JSON License
	Signature string `json:"signature"` // base64 dari RSA-SHA256 signature atas payload
}

// Mode deployment yang didukung.
const (
	ModeSaaS      = "saas"
	ModeOnPremise = "on_premise"
)

// Errors standar untuk verifikasi lisensi.
var (
	ErrExpired        = errors.New("license has expired")
	ErrInvalidSig     = errors.New("invalid license signature")
	ErrMalformed      = errors.New("malformed license file")
	ErrKeyUnavailable = errors.New("RSA key unavailable")
)

// DefaultKeySize adalah ukuran RSA keypair yang di-generate.
const DefaultKeySize = 2048

// GenerateKeyPair membuat RSA keypair baru dan mengembalikan
// PEM-encoded private & public key (format PKCS#8 / PKIX).
func GenerateKeyPair(bits int) (privatePEM, publicPEM []byte, err error) {
	if bits == 0 {
		bits = DefaultKeySize
	}
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("onpremise: generate rsa key: %w", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("onpremise: marshal private key: %w", err)
	}
	privatePEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("onpremise: marshal public key: %w", err)
	}
	publicPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	return privatePEM, publicPEM, nil
}

// SignLicense membuat file .lic (bytes) dari License menggunakan private key PEM.
// Isi payload di-canonical-kan via json.Marshal sebelum di-sign agar
// verifikasi deterministik.
func SignLicense(privateKeyPEM []byte, lic License) ([]byte, error) {
	priv, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	payloadBytes, err := json.Marshal(lic)
	if err != nil {
		return nil, fmt.Errorf("onpremise: marshal license: %w", err)
	}

	digest := sha256.Sum256(payloadBytes)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digest[:])
	if err != nil {
		return nil, fmt.Errorf("onpremise: sign license: %w", err)
	}

	file := licFile{
		Payload:   base64.StdEncoding.EncodeToString(payloadBytes),
		Signature: base64.StdEncoding.EncodeToString(signature),
	}
	out, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("onpremise: marshal lic file: %w", err)
	}
	return out, nil
}

// WriteLicenseFile menulis file .lic ke path.
func WriteLicenseFile(path string, privateKeyPEM []byte, lic License) error {
	data, err := SignLicense(privateKeyPEM, lic)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("onpremise: write license file: %w", err)
	}
	return nil
}

// ReadLicenseFile membaca & memverifikasi file .lic dari disk
// menggunakan public key PEM. Mengembalikan error jika signature invalid
// atau lisensi expired.
func ReadLicenseFile(path, publicKeyPEM string) (*License, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("onpremise: read license file: %w", err)
	}
	return VerifyLicense(publicKeyPEM, data)
}

// VerifyLicense memverifikasi bytes file .lic terhadap public key PEM.
func VerifyLicense(publicKeyPEM string, licData []byte) (*License, error) {
	pub, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	var file licFile
	if err := json.Unmarshal(licData, &file); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMalformed, err)
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(file.Payload)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid payload encoding", ErrMalformed)
	}
	signature, err := base64.StdEncoding.DecodeString(file.Signature)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid signature encoding", ErrMalformed)
	}

	digest := sha256.Sum256(payloadBytes)
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], signature); err != nil {
		return nil, ErrInvalidSig
	}

	var lic License
	if err := json.Unmarshal(payloadBytes, &lic); err != nil {
		return nil, fmt.Errorf("%w: invalid payload json: %v", ErrMalformed, err)
	}

	now := time.Now()
	if lic.ExpiresAt.Before(now) {
		return nil, ErrExpired
	}

	return &lic, nil
}

// IsValidOnPremise mengecek apakah lisensi valid (belum expired).
func (l *License) IsValid() bool {
	return l != nil && time.Now().Before(l.ExpiresAt)
}

// HasModule mengecek apakah modul diizinkan oleh lisensi.
func (l *License) HasModule(slug string) bool {
	if l == nil {
		return false
	}
	for _, m := range l.AllowedModules {
		if m == slug {
			return true
		}
	}
	return false
}

func parsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("%w: invalid private key PEM", ErrKeyUnavailable)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Fallback ke PKCS#1 untuk kompatibilitas key lama.
		if rsaKey, rsaErr := x509.ParsePKCS1PrivateKey(block.Bytes); rsaErr == nil {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("%w: parse private key: %v", ErrKeyUnavailable, err)
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%w: not an RSA private key", ErrKeyUnavailable)
	}
	return rsaKey, nil
}

func parsePublicKey(pemData string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("%w: invalid public key PEM", ErrKeyUnavailable)
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// Fallback ke PKCS#1.
		if rsaKey, rsaErr := x509.ParsePKCS1PublicKey(block.Bytes); rsaErr == nil {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("%w: parse public key: %v", ErrKeyUnavailable, err)
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%w: not an RSA public key", ErrKeyUnavailable)
	}
	return rsaKey, nil
}
