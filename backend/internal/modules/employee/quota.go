package employee

import "errors"

// EmployeeQuotaChecker memeriksa batas kuota jumlah employee.
// Digunakan pada mode deployment on_premise: batas diambil dari
// max_employees pada file .lic. Mode saas → nil (tidak dibatasi).
//
// MaxEmployees() <= 0 berarti tidak ada batas (unlimited).
type EmployeeQuotaChecker interface {
	MaxEmployees() int
}

// ErrQuotaExceeded dikembalikan saat jumlah employee sudah mencapai/melampaui
// batas maksimal yang diizinkan lisensi on-premise.
var ErrQuotaExceeded = errors.New("employee quota exceeded: max employees reached for this license")
