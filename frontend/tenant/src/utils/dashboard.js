// Helper bersama untuk sub-halaman dashboard (dipakai komponen-komponen
// di components/dashboard/ agar tidak duplikasi).

export function localDateStr(d) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

export function formatTime(v) {
  if (!v) return ''
  return String(v).slice(11, 16)
}

export function formatDays(v) {
  const n = Number(v || 0)
  return Number.isInteger(n) ? String(n) : n.toFixed(1)
}

export function fmtScore(v) {
  return (Number(v) || 0).toFixed(1)
}

export function fmtPct(v) {
  return `${(Number(v) || 0).toFixed(1)}%`
}

// Rentang bulan berjalan: { from: 'YYYY-MM-01', to: 'YYYY-MM-DD' }.
export function monthRange(now = new Date()) {
  return { from: `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`, to: localDateStr(now) }
}

// Label minggu (tanggal Senin) format pendek "d/M".
export function shortWeekLabel(s) {
  const d = new Date(String(s).slice(0, 10) + 'T00:00:00')
  if (Number.isNaN(d.getTime())) return s
  return `${d.getDate()}/${d.getMonth() + 1}`
}

// Segmen donut (SVG stroke-dasharray) dari daftar { label, value, color }.
// Dipakai semua chart donut di dashboard (jenis kelamin, status kepegawaian,
// rating KPI/OKR, status sesi, penggunaan cuti).
export function buildDonutSegments(items) {
  const total = items.reduce((s, i) => s + (Number(i.value) || 0), 0)
  if (!total) return []
  const C = 2 * Math.PI * 45
  let acc = 0
  return items.filter(i => (Number(i.value) || 0) > 0).map(i => {
    const frac = (Number(i.value) || 0) / total
    const seg = { ...i, pct: Math.round(frac * 100), dash: `${frac * C} ${C}`, offset: -acc * C }
    acc += frac
    return seg
  })
}
