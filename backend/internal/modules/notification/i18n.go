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
