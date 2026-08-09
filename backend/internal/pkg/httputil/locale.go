package httputil

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// Lang represents a language/locale identifier.
type Lang string

const (
	// LangEN is English (default).
	LangEN Lang = "en"
	// LangID is Bahasa Indonesia.
	LangID Lang = "id"
)

// detectLang reads the Accept-Language header from the request and returns
// the best matching language. Defaults to English.
func detectLang(c *gin.Context) Lang {
	// First, check if language was already set in context by Localize middleware.
	if ctxLang := c.GetString("lang"); ctxLang != "" {
		switch ctxLang {
		case "id", "in":
			return LangID
		default:
			return LangEN
		}
	}

	// Fallback: parse Accept-Language header directly.
	accept := c.GetHeader("Accept-Language")
	if accept == "" {
		return LangEN
	}
	// Accept-Language can be "id", "id-ID", "in", "in-ID", etc.
	lang := strings.Split(accept, ",")[0]
	lang = strings.Split(lang, ";")[0]
	lang = strings.TrimSpace(lang)
	// Normalise: "id", "id-ID", "in", "in-ID" → ID, everything else → EN
	prefix := strings.Split(lang, "-")[0]
	switch strings.ToLower(prefix) {
	case "id", "in":
		return LangID
	default:
		return LangEN
	}
}

// GetLang returns the language for the current request. It checks the Gin
// context first (set by the Localize middleware), and falls back to parsing
// the Accept-Language header if no middleware was used.
func GetLang(c *gin.Context) Lang {
	return detectLang(c)
}

// Translate returns the translated message for the given key and language,
// with optional parameter substitution. Exported so non-HTTP layers (e.g.
// service-level error types) can build message text from the same catalog.
func Translate(lang Lang, key string, params ...string) string {
	return t(lang, key, params...)
}

// t returns the translated message for the given tag and language, with
// optional parameter substitution (%s).
//
// If the language has no entry for the tag, English is used as fallback.
// If English also has no entry, the tag name itself is returned as a
// fallback message.
func t(lang Lang, tag string, params ...string) string {
	msgs, ok := localeMessages[tag]
	if !ok {
		return "Validation failed on '" + tag + "'"
	}

	msg, ok := msgs[lang]
	if !ok {
		msg = msgs[LangEN]
	}

	if len(params) > 0 {
		args := make([]interface{}, len(params))
		for i, p := range params {
			args[i] = p
		}
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// localeMessages is the central message catalog. Each key is a validation tag
// and each value maps language → localised string.
//
// Use %s as a placeholder for tag parameters (e.g. min, max, len).
var localeMessages = map[string]map[Lang]string{
	// --- Meta / special keys ---
	"_title": {
		LangEN: "Validation failed",
		LangID: "Validasi gagal",
	},
	"_unknown": {
		LangEN: "Validation failed on '%s'",
		LangID: "Validasi gagal pada aturan '%s'",
	},

	// --- Built-in validator tags ---
	"required": {
		LangEN: "This field is required",
		LangID: "Field ini wajib diisi",
	},
	"email": {
		LangEN: "Must be a valid email address",
		LangID: "Format email tidak valid",
	},
	"min": {
		LangEN: "Minimum value is %s",
		LangID: "Nilai minimal adalah %s",
	},
	"max": {
		LangEN: "Maximum value is %s",
		LangID: "Nilai maksimal adalah %s",
	},
	"len": {
		LangEN: "Length must be %s",
		LangID: "Panjang harus %s karakter",
	},
	"min_len": {
		LangEN: "Minimum length is %s",
		LangID: "Panjang minimal %s karakter",
	},
	"max_len": {
		LangEN: "Maximum length is %s",
		LangID: "Panjang maksimal %s karakter",
	},
	"oneof": {
		LangEN: "Must be one of: %s",
		LangID: "Nilai harus salah satu dari: %s",
	},
	"numeric": {
		LangEN: "Must be a number",
		LangID: "Harus berupa angka",
	},
	"alpha": {
		LangEN: "Must contain only letters",
		LangID: "Harus berupa huruf",
	},
	"alphanum": {
		LangEN: "Must contain only letters and numbers",
		LangID: "Harus berupa huruf atau angka",
	},
	"url": {
		LangEN: "Must be a valid URL",
		LangID: "Format URL tidak valid",
	},
	"uuid": {
		LangEN: "Must be a valid UUID",
		LangID: "Format UUID tidak valid",
	},
	"datetime": {
		LangEN: "Invalid date format, must be %s",
		LangID: "Format tanggal tidak valid, harus %s",
	},
	"boolean": {
		LangEN: "Must be a boolean value",
		LangID: "Harus berupa boolean",
	},

	// --- Custom Indonesian validators ---
	TagNIK: {
		LangEN: "Invalid NIK format, must be 16 digits",
		LangID: "Format NIK tidak valid, harus 16 digit angka",
	},
	TagNPWP: {
		LangEN: "Invalid NPWP format, must be 15-16 digits",
		LangID: "Format NPWP tidak valid, harus 15-16 digit angka",
	},
	TagPhoneID: {
		LangEN: "Invalid Indonesian phone number format",
		LangID: "Format nomor telepon Indonesia tidak valid",
	},
	TagPhoneIndonesia: {
		LangEN: "Invalid Indonesian phone number format",
		LangID: "Format nomor telepon Indonesia tidak valid",
	},
	TagPostalCode: {
		LangEN: "Invalid postal code, must be 5 digits",
		LangID: "Format kode pos tidak valid, harus 5 digit angka",
	},
	TagKodePos: {
		LangEN: "Invalid postal code, must be 5 digits",
		LangID: "Format kode pos tidak valid, harus 5 digit angka",
	},
	TagDateID: {
		LangEN: "Invalid date format, must be YYYY-MM-DD",
		LangID: "Format tanggal tidak valid, harus YYYY-MM-DD",
	},
	TagKK: {
		LangEN: "Invalid KK format, must be 16 digits",
		LangID: "Format KK tidak valid, harus 16 digit angka",
	},
	TagNoRekening: {
		LangEN: "Invalid bank account number, must be 8-20 digits",
		LangID: "Format nomor rekening tidak valid, harus 8-20 digit angka",
	},
	TagPassport: {
		LangEN: "Invalid passport number, must be 1 letter followed by 8 digits",
		LangID: "Format nomor passport tidak valid, harus 1 huruf diikuti 8 angka",
	},
	TagSIM: {
		LangEN: "Invalid SIM number, must be 12 digits",
		LangID: "Format SIM tidak valid, harus 12 digit angka",
	},
	TagNPWPFormatted: {
		LangEN: "Invalid NPWP format, must be XX.XXX.XXX.X-XXX.XXX",
		LangID: "Format NPWP tidak valid, harus XX.XXX.XXX.X-XXX.XXX",
	},

	// --- Error / status code messages ---
	"error.not_found": {
		LangEN: "Not found",
		LangID: "Tidak ditemukan",
	},
	"error.internal": {
		LangEN: "Internal server error",
		LangID: "Terjadi kesalahan pada server",
	},
	"error.unauthorized": {
		LangEN: "Unauthorized",
		LangID: "Tidak terotorisasi",
	},
	"error.forbidden": {
		LangEN: "Forbidden",
		LangID: "Akses ditolak",
	},
	"error.conflict": {
		LangEN: "Conflict",
		LangID: "Konflik",
	},
	"error.bad_request": {
		LangEN: "Bad request",
		LangID: "Permintaan tidak valid",
	},
	"error.validation": {
		LangEN: "Validation failed",
		LangID: "Validasi gagal",
	},
	"error.user_not_authenticated": {
		LangEN: "User not authenticated",
		LangID: "Pengguna belum terautentikasi",
	},
	"error.invalid_user_context": {
		LangEN: "Invalid user id in context",
		LangID: "ID pengguna tidak valid di konteks",
	},
	"error.invalid_id": {
		LangEN: "Invalid id",
		LangID: "ID tidak valid",
	},

	// --- Success action messages ---
	"success.created": {
		LangEN: "Created successfully",
		LangID: "Berhasil dibuat",
	},
	"success.updated": {
		LangEN: "Updated successfully",
		LangID: "Berhasil diperbarui",
	},
	"success.deleted": {
		LangEN: "Deleted successfully",
		LangID: "Berhasil dihapus",
	},
	"success.saved": {
		LangEN: "Saved successfully",
		LangID: "Berhasil disimpan",
	},
	"success.approved": {
		LangEN: "Approved successfully",
		LangID: "Berhasil disetujui",
	},
	"success.rejected": {
		LangEN: "Rejected successfully",
		LangID: "Berhasil ditolak",
	},
	"success.executed": {
		LangEN: "Executed successfully",
		LangID: "Berhasil dieksekusi",
	},
	"success.actuals_updated": {
		LangEN: "Actuals updated successfully",
		LangID: "Nilai aktual berhasil diperbarui",
	},
	"success.cancelled": {
		LangEN: "Cancelled successfully",
		LangID: "Berhasil dibatalkan",
	},
	"success.restored": {
		LangEN: "Restored successfully",
		LangID: "Berhasil dipulihkan",
	},
	"success.suspended": {
		LangEN: "Suspended successfully",
		LangID: "Berhasil dinonaktifkan",
	},
	"success.activated": {
		LangEN: "Activated successfully",
		LangID: "Berhasil diaktifkan",
	},
	"success.terminated": {
		LangEN: "Terminated successfully",
		LangID: "Berhasil dihentikan",
	},
	"success.backed_up": {
		LangEN: "Backup completed successfully",
		LangID: "Backup berhasil",
	},
	"success.restored_data": {
		LangEN: "Restore completed successfully",
		LangID: "Pemulihan berhasil",
	},

	// --- Module-specific messages ---
	"company.created": {
		LangEN: "Company created successfully",
		LangID: "Perusahaan berhasil dibuat",
	},
	"company.updated": {
		LangEN: "Company updated successfully",
		LangID: "Perusahaan berhasil diperbarui",
	},
	"company.deleted": {
		LangEN: "Company deleted successfully",
		LangID: "Perusahaan berhasil dihapus",
	},
	"company.suspended": {
		LangEN: "Company suspended successfully — tenant connection deactivated",
		LangID: "Perusahaan berhasil dinonaktifkan — koneksi tenant dinonaktifkan",
	},
	"company.activated": {
		LangEN: "Company activated successfully — tenant connection reactivated",
		LangID: "Perusahaan berhasil diaktifkan — koneksi tenant diaktifkan kembali",
	},
	"company.terminated": {
		LangEN: "Company terminated — tenant database dropped and connection removed",
		LangID: "Perusahaan dihentikan — database tenant dihapus dan koneksi dicabut",
	},

	"user.created": {
		LangEN: "User created successfully",
		LangID: "Pengguna berhasil dibuat",
	},
	"user.updated": {
		LangEN: "User updated successfully",
		LangID: "Pengguna berhasil diperbarui",
	},
	"user.deleted": {
		LangEN: "User deleted successfully",
		LangID: "Pengguna berhasil dihapus",
	},
	"password.updated": {
		LangEN: "Password updated successfully",
		LangID: "Kata sandi berhasil diperbarui",
	},

	"license.created": {
		LangEN: "License created successfully",
		LangID: "Lisensi berhasil dibuat",
	},
	"license.updated": {
		LangEN: "License updated successfully",
		LangID: "Lisensi berhasil diperbarui",
	},

	"module.created": {
		LangEN: "Module created successfully",
		LangID: "Modul berhasil dibuat",
	},
	"module.updated": {
		LangEN: "Module updated successfully",
		LangID: "Modul berhasil diperbarui",
	},
	"module.deleted": {
		LangEN: "Module deleted successfully",
		LangID: "Modul berhasil dihapus",
	},

	"approval.flow.created": {
		LangEN: "Approval flow created",
		LangID: "Alur persetujuan berhasil dibuat",
	},
	"approval.flow.updated": {
		LangEN: "Approval flow updated",
		LangID: "Alur persetujuan berhasil diperbarui",
	},
	"approval.flow.deleted": {
		LangEN: "Approval flow deleted",
		LangID: "Alur persetujuan berhasil dihapus",
	},
	"approval.step.created": {
		LangEN: "Approval step created",
		LangID: "Langkah persetujuan berhasil dibuat",
	},
	"approval.step.updated": {
		LangEN: "Approval step updated",
		LangID: "Langkah persetujuan berhasil diperbarui",
	},
	"approval.step.deleted": {
		LangEN: "Approval step deleted",
		LangID: "Langkah persetujuan berhasil dihapus",
	},
	"approval.instance.created": {
		LangEN: "Approval instance created",
		LangID: "Instance persetujuan berhasil dibuat",
	},
	"approval.instance.cancelled": {
		LangEN: "Approval instance cancelled",
		LangID: "Instance persetujuan berhasil dibatalkan",
	},
	"approval.action.submitted": {
		LangEN: "Approval action submitted",
		LangID: "Aksi persetujuan berhasil dikirim",
	},
	"approval.module_param_required": {
		LangEN: "module query parameter is required",
		LangID: "Parameter query module wajib diisi",
	},
	"approval.no_submitter": {
		LangEN: "Cannot resolve supervisor: instance has no submitter",
		LangID: "Tidak dapat menentukan supervisor: instance tidak memiliki pengaju",
	},
	"approval.submitter_org_unresolved": {
		LangEN: "Submitter's organization could not be resolved",
		LangID: "Organisasi pengaju tidak dapat ditentukan",
	},
	"approval.no_org_levels": {
		LangEN: "No organization %s level(s) above the submitter for step \"%s\"",
		LangID: "Tidak ada organisasi %s level di atas pengaju untuk langkah \"%s\"",
	},
	"approval.no_supervisor_vacant": {
		LangEN: "No supervisor found for step \"%s\": every organization from the submitter up to the top of the hierarchy is vacant",
		LangID: "Supervisor tidak ditemukan untuk langkah \"%s\": seluruh organisasi dari pengaju hingga puncak hierarki kosong",
	},
	"approval.no_assignees": {
		LangEN: "No assignees resolved for step: %s",
		LangID: "Tidak ada penyetuju yang dapat ditentukan untuk langkah: %s",
	},
	"approval.flow_inactive": {
		LangEN: "Flow is not active",
		LangID: "Alur persetujuan tidak aktif",
	},
	"approval.flow_no_steps": {
		LangEN: "Flow has no steps configured",
		LangID: "Alur persetujuan belum memiliki langkah",
	},
	"approval.pending_instance_exists": {
		LangEN: "Document already has a pending approval instance",
		LangID: "Dokumen sudah memiliki instance persetujuan yang berjalan",
	},
	"approval.zero_approvers": {
		LangEN: "Approval step \"%s\" (order %s) resolved to zero approvers — check the flow's approver_type/organization/role configuration for this step",
		LangID: "Langkah persetujuan \"%s\" (urutan %s) tidak memiliki penyetuju — periksa konfigurasi approver_type/organization/role pada langkah ini",
	},

	"employee.created": {
		LangEN: "Employee created successfully",
		LangID: "Karyawan berhasil dibuat",
	},
	"employee.updated": {
		LangEN: "Employee updated successfully",
		LangID: "Karyawan berhasil diperbarui",
	},
	"employee.deleted": {
		LangEN: "Employee deleted",
		LangID: "Karyawan berhasil dihapus",
	},

	"organization.created": {
		LangEN: "Organization created successfully",
		LangID: "Organisasi berhasil dibuat",
	},
	"organization.updated": {
		LangEN: "Organization updated successfully",
		LangID: "Organisasi berhasil diperbarui",
	},
	"organization.deleted": {
		LangEN: "Organization deleted",
		LangID: "Organisasi berhasil dihapus",
	},

	"shift.deleted": {
		LangEN: "Shift deleted",
		LangID: "Shift berhasil dihapus",
	},
	"employee_shift.deleted": {
		LangEN: "Employee shift deleted",
		LangID: "Shift karyawan berhasil dihapus",
	},
	"location.deleted": {
		LangEN: "Location deleted",
		LangID: "Lokasi berhasil dihapus",
	},
	"overtime.deleted": {
		LangEN: "Overtime request deleted",
		LangID: "Lembur berhasil dihapus",
	},
	"exempt_position.deleted": {
		LangEN: "Exempt position deleted",
		LangID: "Posisi exempt berhasil dihapus",
	},

	"payroll.component.created": {
		LangEN: "Salary component created",
		LangID: "Komponen gaji berhasil dibuat",
	},
	"payroll.component.deleted": {
		LangEN: "Salary component deleted",
		LangID: "Komponen gaji berhasil dihapus",
	},
	"payroll.profile.deleted": {
		LangEN: "Employee payroll profile deleted",
		LangID: "Profil penggajian berhasil dihapus",
	},

	"leave.type.deleted": {
		LangEN: "Leave type deleted",
		LangID: "Tipe cuti berhasil dihapus",
	},
	"leave.policy.deleted": {
		LangEN: "Accrual policy deleted",
		LangID: "Kebijakan akrual berhasil dihapus",
	},
	"leave.reason.deleted": {
		LangEN: "Leave reason deleted",
		LangID: "Alasan cuti berhasil dihapus",
	},
	"leave.request.deleted": {
		LangEN: "Leave request deleted",
		LangID: "Permohonan cuti berhasil dihapus",
	},

	"training.category.deleted": {
		LangEN: "Training category deleted",
		LangID: "Kategori pelatihan berhasil dihapus",
	},
	"training.course.deleted": {
		LangEN: "Training course deleted",
		LangID: "Kursus pelatihan berhasil dihapus",
	},
	"training.session.deleted": {
		LangEN: "Training session deleted",
		LangID: "Sesi pelatihan berhasil dihapus",
	},
	"training.participant.deleted": {
		LangEN: "Training participant deleted",
		LangID: "Peserta pelatihan berhasil dihapus",
	},
	"training.material.deleted": {
		LangEN: "Training material deleted",
		LangID: "Materi pelatihan berhasil dihapus",
	},
	"training.evaluation.deleted": {
		LangEN: "Training evaluation deleted",
		LangID: "Evaluasi pelatihan berhasil dihapus",
	},
	"training.certificate.deleted": {
		LangEN: "Training certificate deleted",
		LangID: "Sertifikat pelatihan berhasil dihapus",
	},

	"career.talent_map.deleted": {
		LangEN: "Talent map deleted",
		LangID: "Peta bakat berhasil dihapus",
	},
	"career.path.deleted": {
		LangEN: "Career path deleted",
		LangID: "Jalur karir berhasil dihapus",
	},
	"career.plan.deleted": {
		LangEN: "Succession plan deleted",
		LangID: "Rencana suksesi berhasil dihapus",
	},

	"rbac.role.deleted": {
		LangEN: "Role deleted",
		LangID: "Peran berhasil dihapus",
	},
	"rbac.permission.deleted": {
		LangEN: "Permission deleted",
		LangID: "Izin berhasil dihapus",
	},
	"rbac.permission.assigned": {
		LangEN: "Permission assigned to role",
		LangID: "Izin berhasil ditetapkan ke peran",
	},
	"rbac.permission.revoked": {
		LangEN: "Permission revoked from role",
		LangID: "Izin berhasil dicabut dari peran",
	},
	"rbac.permissions_updated": {
		LangEN: "Permissions updated successfully",
		LangID: "Permissions berhasil diperbarui",
	},
	"rbac.user_roles_updated": {
		LangEN: "User roles updated successfully",
		LangID: "Roles user berhasil diperbarui",
	},
	"rbac.duplicate_role": {
		LangEN: "Role '%s' already exists",
		LangID: "Role '%s' sudah ada",
	},
	"rbac.system_role": {
		LangEN: "Role '%s' is a system role and cannot be modified",
		LangID: "Role '%s' adalah role sistem dan tidak dapat diubah",
	},

	"address.deleted": {
		LangEN: "Address deleted",
		LangID: "Alamat berhasil dihapus",
	},
	"contact.deleted": {
		LangEN: "Emergency contact deleted",
		LangID: "Kontak darurat berhasil dihapus",
	},
	"family.deleted": {
		LangEN: "Family member deleted",
		LangID: "Anggota keluarga berhasil dihapus",
	},
	"education.deleted": {
		LangEN: "Education record deleted",
		LangID: "Riwayat pendidikan berhasil dihapus",
	},
	"experience.deleted": {
		LangEN: "Work experience deleted",
		LangID: "Pengalaman kerja berhasil dihapus",
	},
	"document.deleted": {
		LangEN: "Document deleted",
		LangID: "Dokumen berhasil dihapus",
	},
	"insurance.deleted": {
		LangEN: "Insurance record deleted",
		LangID: "Asuransi berhasil dihapus",
	},
	"employment.deleted": {
		LangEN: "Employment record deleted",
		LangID: "Riwayat pekerjaan berhasil dihapus",
	},

	// --- Setting module messages ---
	"setting.duplicate_code": {
		LangEN: "Code '%s' is already in use",
		LangID: "Kode '%s' sudah digunakan",
	},

	// --- Additional module messages ---
	"training.session_id_required": {
		LangEN: "session_id is required",
		LangID: "session_id wajib diisi",
	},
	"attendance.query_params_required": {
		LangEN: "employee_id and work_date query parameters are required",
		LangID: "Parameter query employee_id dan work_date wajib diisi",
	},
	"payroll.instance_param_required": {
		LangEN: "instance_id query parameter is required",
		LangID: "Parameter query instance_id wajib diisi",
	},
	"performance.template_param_required": {
		LangEN: "template_id query parameter is required",
		LangID: "Parameter query template_id wajib diisi",
	},
	"performance.employee_id_required": {
		LangEN: "employee_id is required",
		LangID: "Parameter employee_id wajib diisi",
	},
	"performance.manager_id_required": {
		LangEN: "manager_id is required",
		LangID: "Parameter manager_id wajib diisi",
	},
	"upload.file_required": {
		LangEN: "Please upload a file",
		LangID: "Mohon unggah sebuah file",
	},
	"upload.type_not_allowed": {
		LangEN: "File type not allowed",
		LangID: "Tipe file tidak diizinkan",
	},
	"upload.file_too_large": {
		LangEN: "File too large (max 10 MB)",
		LangID: "File terlalu besar (maks 10 MB)",
	},
	"employee.photo_required": {
		LangEN: "Photo file is required",
		LangID: "File foto wajib diunggah",
	},
	"employee.image_type_only": {
		LangEN: "Only image files (jpg, jpeg, png, gif, webp) are allowed",
		LangID: "Hanya file gambar yang diizinkan (jpg, jpeg, png, gif, webp)",
	},
	"employee.file_max_2mb": {
		LangEN: "File size must be less than 2MB",
		LangID: "Ukuran file maksimal 2MB",
	},
	"employee.name_required": {
		LangEN: "Name is required",
		LangID: "Nama wajib diisi",
	},
	"employee.document_required": {
		LangEN: "Document file is required",
		LangID: "File dokumen wajib diunggah",
	},
	"employee.file_type_not_allowed": {
		LangEN: "File type not allowed. Allowed: PDF, DOC, DOCX, XLS, XLSX, JPG, PNG, GIF, TXT",
		LangID: "Tipe file tidak diizinkan. Diizinkan: PDF, DOC, DOCX, XLS, XLSX, JPG, PNG, GIF, TXT",
	},
	"employee.file_type_not_allowed_short": {
		LangEN: "File type not allowed",
		LangID: "Tipe file tidak diizinkan",
	},
	"employee.file_max_10mb": {
		LangEN: "File size must be less than 10MB",
		LangID: "Ukuran file maksimal 10MB",
	},
	"employee.invalid_id": {
		LangEN: "Invalid employee ID",
		LangID: "ID karyawan tidak valid",
	},

	// --- Package Management messages ---
	"package.created": {
		LangEN: "Package created successfully",
		LangID: "Paket berhasil dibuat",
	},
	"package.updated": {
		LangEN: "Package updated successfully",
		LangID: "Paket berhasil diperbarui",
	},
	"package.deleted": {
		LangEN: "Package deleted successfully",
		LangID: "Paket berhasil dihapus",
	},
	"package.published": {
		LangEN: "Package published successfully",
		LangID: "Paket berhasil dipublikasikan",
	},
	"package.unpublished": {
		LangEN: "Package unpublished successfully",
		LangID: "Paket berhasil ditarik",
	},
	"package.subscribed": {
		LangEN: "Successfully subscribed to package",
		LangID: "Berhasil berlangganan paket",
	},
	"package.unsubscribed": {
		LangEN: "Successfully unsubscribed from package",
		LangID: "Berhasil berhenti berlangganan paket",
	},
}
