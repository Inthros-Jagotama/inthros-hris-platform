package calculator

import "sort"

// VariableMeta mendeskripsikan sebuah variabel built-in yang dikenal formula engine.
type VariableMeta struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// Registry menyimpan daftar variabel built-in yang boleh dipakai di formula
// payroll. Variabel selain ini dianggap referensi ke salary component lain
// (diresolusi oleh payroll run saat evaluasi).
type Registry struct {
	variables map[string]VariableMeta
}

// NewRegistry membuat Registry kosong.
func NewRegistry() *Registry {
	return &Registry{variables: map[string]VariableMeta{}}
}

// Register menambahkan sebuah variabel built-in. Nama di-upper-case otomatis.
func (r *Registry) Register(meta VariableMeta) {
	meta.Name = normalizeVarName(meta.Name)
	r.variables[meta.Name] = meta
}

// IsBuiltIn mengembalikan true jika nama variabel (case-insensitive) terdaftar.
func (r *Registry) IsBuiltIn(name string) bool {
	_, ok := r.variables[normalizeVarName(name)]
	return ok
}

// Get mengambil metadata variabel built-in.
func (r *Registry) Get(name string) (VariableMeta, bool) {
	meta, ok := r.variables[normalizeVarName(name)]
	return meta, ok
}

// All mengembalikan semua variabel built-in, diurutkan berdasarkan nama.
func (r *Registry) All() []VariableMeta {
	names := make([]string, 0, len(r.variables))
	for name := range r.variables {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]VariableMeta, 0, len(names))
	for _, name := range names {
		out = append(out, r.variables[name])
	}
	return out
}

// normalizeVarName menyeragamkan nama variabel menjadi UPPER_SNAKE.
func normalizeVarName(name string) string {
	var out []byte
	for i := 0; i < len(name); i++ {
		c := name[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

// DefaultRegistry mengembalikan registry dengan variabel built-in standar
// payroll (lihat docs/payroll/02-formula-engine.md §Variable registry).
func DefaultRegistry() *Registry {
	r := NewRegistry()
	inputs := []VariableMeta{
		{Name: "BASIC", Category: "STRUCTURE", Description: "Basic salary komponen BASIC"},
		{Name: "GROSS", Category: "TOTAL", Description: "Total penghasilan kotor (total earnings)"},
		{Name: "BPJS_WAGE", Category: "STATUTORY", Description: "Dasar upah perhitungan BPJS"},
		{Name: "TAXABLE_INCOME", Category: "STATUTORY", Description: "Penghasilan kena pajak (PPh21)"},
		{Name: "WORKING_DAYS", Category: "WORKFORCE", Description: "Jumlah hari kerja dalam periode"},
		{Name: "WORKED_DAYS", Category: "WORKFORCE", Description: "Jumlah hari kerja yang dihadiri"},
		{Name: "ABSENCE_DAYS", Category: "WORKFORCE", Description: "Jumlah hari absen/alpha"},
		{Name: "UNPAID_LEAVE_DAYS", Category: "WORKFORCE", Description: "Jumlah hari cuti tidak dibayar"},
		{Name: "OVERTIME_HOURS", Category: "WORKFORCE", Description: "Total jam lembur"},
		{Name: "OVERTIME_RATE", Category: "WORKFORCE", Description: "Tarif per jam lembur"},
		{Name: "TOTAL_EARNINGS", Category: "TOTAL", Description: "Total seluruh komponen earning"},
		{Name: "TOTAL_DEDUCTIONS", Category: "TOTAL", Description: "Total seluruh komponen deduction"},
		{Name: "TOTAL_EMPLOYEE_DEDUCTION", Category: "TOTAL", Description: "Total potongan karyawan (BPJS + PPh21 + lainnya)"},
		{Name: "TOTAL_EMPLOYER_CONTRIBUTION", Category: "TOTAL", Description: "Total iuran/contribution perusahaan"},
		{Name: "NET_SALARY", Category: "TOTAL", Description: "Gaji bersih setelah potongan"},
		{Name: "EMPLOYER_TOTAL_COST", Category: "TOTAL", Description: "Total biaya perusahaan (gross + employer contribution)"},
	}
	for _, meta := range inputs {
		r.Register(meta)
	}
	return r
}
