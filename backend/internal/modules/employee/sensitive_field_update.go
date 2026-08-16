package employee

import (
	"github.com/inthros/hris-platform/internal/pkg/crypto"
	"github.com/inthros/hris-platform/internal/pkg/mask"
)

// =============================================================================
// Masked-echo guard untuk update field sensitif.
//
// GET employee mengembalikan nilai ter-mask (mis. "************5678") ke caller
// yang tidak punya permission "view_*". Form edit di frontend memuat response
// itu apa adanya lalu mengirim kembali SELURUH objek saat save — termasuk field
// sensitif yang tidak pernah disentuh user. Tanpa guard, string mask tersebut
// akan ditulis ke DB (dan ikut dienkripsi bila toggle aktif), menghancurkan PII
// asli tanpa jejak.
//
// Guard di bawah membandingkan nilai yang masuk dengan hasil mask dari nilai
// tersimpan saat ini (didekripsi lebih dulu bila perlu). Jika identik, request
// dianggap "tidak mengubah apa-apa" dan penulisan di-skip.
//
// API sengaja jadi trust boundary di sini: perbaikan di frontend saja tidak
// cukup, karena klien mana pun bisa mengirim PUT yang sama.
// =============================================================================

// isMaskedEcho melaporkan apakah incoming hanyalah pantulan bentuk ter-mask
// dari stored (nilai tersimpan, boleh ciphertext), bukan edit sungguhan.
func isMaskedEcho(stored, incoming string) bool {
	if stored == "" || incoming == "" {
		return false
	}
	plain := crypto.DecryptIfLooksEncrypted(stored)
	if incoming == plain {
		// Nilai yang sama persis dengan plaintext tersimpan — bukan mask.
		return false
	}
	return incoming == mask.PartialMask(plain)
}

// applyIfNotMasked menyalin incoming ke target kecuali incoming hanyalah
// pantulan mask dari nilai yang sudah tersimpan di target.
// Dipakai untuk field sensitif bertipe string (non-pointer) pada model.
func applyIfNotMasked(incoming *string, target *string) {
	if incoming == nil {
		return
	}
	if isMaskedEcho(*target, *incoming) {
		return
	}
	*target = *incoming
}

// applyPtrIfNotMasked adalah varian applyIfNotMasked untuk field sensitif
// bertipe *string pada model (nullable di DB).
func applyPtrIfNotMasked(incoming *string, target **string) {
	if incoming == nil {
		return
	}
	stored := ""
	if *target != nil {
		stored = **target
	}
	if isMaskedEcho(stored, *incoming) {
		return
	}
	*target = incoming
}
