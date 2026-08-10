package notification

import (
	"fmt"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// catalogEntry holds the EN/ID title + body templates for one notification
// type. Body templates use %s placeholders, filled positionally from the
// notification's stored Params at render time (see translate below).
type catalogEntry struct {
	title map[httputil.Lang]string
	body  map[httputil.Lang]string
}

// catalog is the bilingual message table for every notification type this
// codebase currently produces (leave, attendance overtime/correction, and
// the central Approval module's own "you have a pending task" pushes).
// Consumer modules pass only Type + Params when calling Notify — the actual
// title/body text is rendered here, in the recipient's own language,
// determined from THEIR GET /notifications request (Accept-Language), not
// whatever language the action that triggered the notification happened to
// run in.
var catalog = map[string]catalogEntry{
	"LEAVE_APPROVED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Leave Request Approved",
			httputil.LangID: "Permohonan Cuti Disetujui",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your leave request has been approved.",
			httputil.LangID: "Permohonan cuti Anda telah disetujui.",
		},
	},
	"LEAVE_REJECTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Leave Request Rejected",
			httputil.LangID: "Permohonan Cuti Ditolak",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your leave request has been rejected.",
			httputil.LangID: "Permohonan cuti Anda telah ditolak.",
		},
	},
	"LEAVE_CANCELLED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Leave Request Cancelled",
			httputil.LangID: "Permohonan Cuti Dibatalkan",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your leave request has been cancelled.",
			httputil.LangID: "Permohonan cuti Anda telah dibatalkan.",
		},
	},
	"MOVEMENT_SUBMITTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Career Movement Submitted",
			httputil.LangID: "Pergerakan Karier Diajukan",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your career movement has been submitted for approval.",
			httputil.LangID: "Pergerakan karier Anda telah diajukan untuk persetujuan.",
		},
	},
	"MOVEMENT_APPROVED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Career Movement Approved",
			httputil.LangID: "Pergerakan Karier Disetujui",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your career movement has been approved. HR can now execute it.",
			httputil.LangID: "Pergerakan karier Anda telah disetujui. HR dapat mengeksekusinya.",
		},
	},
	"MOVEMENT_REJECTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Career Movement Rejected",
			httputil.LangID: "Pergerakan Karier Ditolak",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your career movement has been rejected.",
			httputil.LangID: "Pergerakan karier Anda telah ditolak.",
		},
	},
	"MOVEMENT_EXECUTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Career Movement Executed",
			httputil.LangID: "Pergerakan Karier Dieksekusi",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your career movement has been executed.",
			httputil.LangID: "Pergerakan karier Anda telah dieksekusi.",
		},
	},
	"CONTRACT_EXPIRING": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Contract Expiring Soon",
			httputil.LangID: "Kontrak Akan Berakhir",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Contract %s expires on %s. Please prepare the renewal.",
			httputil.LangID: "Kontrak %s berakhir pada %s. Silakan siapkan perpanjangan.",
		},
	},
	"CONTRACT_EXPIRED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Contract Expired",
			httputil.LangID: "Kontrak Berakhir",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Contract %s has expired.",
			httputil.LangID: "Kontrak %s telah berakhir.",
		},
	},
	"OVERTIME_APPROVED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Overtime Request Approved",
			httputil.LangID: "Pengajuan Lembur Disetujui",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your overtime request has been approved.",
			httputil.LangID: "Pengajuan lembur Anda telah disetujui.",
		},
	},
	"OVERTIME_REJECTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Overtime Request Rejected",
			httputil.LangID: "Pengajuan Lembur Ditolak",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your overtime request has been rejected.",
			httputil.LangID: "Pengajuan lembur Anda telah ditolak.",
		},
	},
	"OVERTIME_ASSIGNED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Overtime Assigned",
			httputil.LangID: "Penugasan Lembur",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "You have been assigned an overtime duty on %s. Please fill in the actual hours.",
			httputil.LangID: "Anda ditugaskan lembur pada %s. Silakan isi jam aktual lembur Anda.",
		},
	},
	"OVERTIME_ACTUAL_APPROVED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Overtime Actual Approved",
			httputil.LangID: "Isian Aktual Lembur Disetujui",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your overtime actual entry has been approved.",
			httputil.LangID: "Isian aktual lembur Anda telah disetujui.",
		},
	},
	"OVERTIME_ACTUAL_REJECTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Overtime Actual Rejected",
			httputil.LangID: "Isian Aktual Lembur Ditolak",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your overtime actual entry has been rejected.",
			httputil.LangID: "Isian aktual lembur Anda telah ditolak.",
		},
	},
	"CORRECTION_APPROVED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Attendance Correction Approved",
			httputil.LangID: "Koreksi Kehadiran Disetujui",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your attendance correction request has been approved.",
			httputil.LangID: "Pengajuan koreksi kehadiran Anda telah disetujui.",
		},
	},
	"CORRECTION_REJECTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Attendance Correction Rejected",
			httputil.LangID: "Koreksi Kehadiran Ditolak",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your attendance correction request has been rejected.",
			httputil.LangID: "Pengajuan koreksi kehadiran Anda telah ditolak.",
		},
	},
	"KPI_TEMPLATE_CREATED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "New KPI Template",
			httputil.LangID: "Template KPI Baru",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "A new KPI template \"%s\" is available for your organization.",
			httputil.LangID: "Template KPI baru \"%s\" tersedia untuk organisasi Anda.",
		},
	},
	"OKR_TEMPLATE_CREATED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "New OKR Template",
			httputil.LangID: "Template OKR Baru",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "A new OKR template \"%s\" is available for your organization.",
			httputil.LangID: "Template OKR baru \"%s\" tersedia untuk organisasi Anda.",
		},
	},
	"KPI_TARGET_APPROVED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "KPI Target Approved",
			httputil.LangID: "Target KPI Disetujui",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your KPI target has been approved.",
			httputil.LangID: "Target KPI Anda telah disetujui.",
		},
	},
	"KPI_TARGET_REJECTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "KPI Target Rejected",
			httputil.LangID: "Target KPI Ditolak",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your KPI target has been rejected. Please revise and resubmit.",
			httputil.LangID: "Target KPI Anda telah ditolak. Silakan perbaiki dan ajukan kembali.",
		},
	},
	"KPI_REALIZATION_APPROVED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "KPI Realization Approved",
			httputil.LangID: "Realisasi KPI Disetujui",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your KPI realization has been approved.",
			httputil.LangID: "Realisasi KPI Anda telah disetujui.",
		},
	},
	"KPI_REALIZATION_REJECTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "KPI Realization Rejected",
			httputil.LangID: "Realisasi KPI Ditolak",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your KPI realization has been rejected. Please revise and resubmit.",
			httputil.LangID: "Realisasi KPI Anda telah ditolak. Silakan perbaiki dan ajukan kembali.",
		},
	},
	"OKR_KEY_RESULT_APPROVED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "OKR Key Results Approved",
			httputil.LangID: "Key Results OKR Disetujui",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your OKR key results have been approved.",
			httputil.LangID: "Key results OKR Anda telah disetujui.",
		},
	},
	"OKR_KEY_RESULT_REJECTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "OKR Key Results Rejected",
			httputil.LangID: "Key Results OKR Ditolak",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your OKR key results have been rejected. Please revise and resubmit.",
			httputil.LangID: "Key results OKR Anda telah ditolak. Silakan perbaiki dan ajukan kembali.",
		},
	},
	"OKR_ASSESSMENT_APPROVED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "OKR Assessment Approved",
			httputil.LangID: "Penilaian OKR Disetujui",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your OKR assessment has been approved.",
			httputil.LangID: "Penilaian OKR Anda telah disetujui.",
		},
	},
	"OKR_ASSESSMENT_REJECTED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "OKR Assessment Rejected",
			httputil.LangID: "Penilaian OKR Ditolak",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "Your OKR assessment has been rejected. Please revise and resubmit.",
			httputil.LangID: "Penilaian OKR Anda telah ditolak. Silakan perbaiki dan ajukan kembali.",
		},
	},
	"APPROVAL_TASK_ASSIGNED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "Approval Needed",
			httputil.LangID: "Perlu Persetujuan",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "A %s request needs your approval (%s).",
			httputil.LangID: "Permohonan %s memerlukan persetujuan Anda (%s).",
		},
	},
	"APPROVAL_WATCHER_ASSIGNED": {
		title: map[httputil.Lang]string{
			httputil.LangEN: "You're Watching an Approval",
			httputil.LangID: "Anda Memantau Sebuah Persetujuan",
		},
		body: map[httputil.Lang]string{
			httputil.LangEN: "A %s request has reached a step you're watching (%s).",
			httputil.LangID: "Permohonan %s telah mencapai langkah yang Anda pantau (%s).",
		},
	},
}

// translate renders a notification's title/body in the given language,
// filling body %s placeholders from params. Falls back to defaultTitle/
// defaultBody (the English text computed once at Notify()-time and stored
// alongside Type/Params) when the type has no catalog entry — a module
// using a notifType this catalog hasn't been taught about yet degrades to
// its own English text instead of showing blank/garbled content.
func translate(l httputil.Lang, notifType string, params []string, defaultTitle, defaultBody string) (string, string) {
	entry, ok := catalog[notifType]
	if !ok {
		return defaultTitle, defaultBody
	}

	title := entry.title[l]
	if title == "" {
		title = entry.title[httputil.LangEN]
	}

	bodyTmpl := entry.body[l]
	if bodyTmpl == "" {
		bodyTmpl = entry.body[httputil.LangEN]
	}
	body := bodyTmpl
	if len(params) > 0 {
		args := make([]interface{}, len(params))
		for i, p := range params {
			args[i] = p
		}
		body = fmt.Sprintf(bodyTmpl, args...)
	}
	return title, body
}
