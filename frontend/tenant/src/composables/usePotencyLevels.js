import { ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

const apiBase = '/api/v1/tenant/job-management/potency-competencies'

/**
 * usePotencyLevels — Logika bersama simpan/hapus/hydrate untuk card level potensi
 * (PsychologicalPotencyCard & SkillPotencyCard).
 *
 * @param {Object} opts
 * @param {import('vue').Ref<string>} opts.orgId — Ref id organisasi (bisa computed(() => props.orgId))
 * @param {import('vue').Ref<Array>} opts.rows — Ref baris yang disimpan ({competency_id, job_management_value_id, recordId, levelOptions})
 * @param {Function} [opts.afterDelete] — Perilaku khusus setelah record terhapus (opsional; psych: lepas tipe dari pilihan)
 * @param {Function} [opts.onSaved] — Callback setelah simpan/hapus berhasil (mis. emit('saved'))
 * @param {string} [opts.matchBy] — Cara mencocokkan record saat hydrate: 'value' (default, baris tetap tanpa
 *   competency_id dicocokkan via job_management_value_id) atau 'competency' (baris berbasis kompetensi
 *   dicocokkan via competency_id — dipakai card Kompetensi Teknis).
 * @param {string} [opts.descriptionField] — Field option yang ditampilkan di kolom deskripsi level:
 *   'descriptions' (default) atau 'note' (dipakai card Kompetensi Teknis).
 * @returns
 */
export function usePotencyLevels({ orgId, rows, afterDelete, onSaved, matchBy = 'value', descriptionField = 'descriptions' }) {
  const { t } = useI18n()
  const toast = useToast()

  const savingCard = ref(false)
  const errorMsg = ref('')
  const deleteVisible = ref(false)
  const deleting = ref(false)
  const deleteError = ref('')
  const deleteTarget = ref(null)

  // Raw record potency dari API — dasar hydration (reaktif agar bisa dipakai luar)
  const records = ref([])

  // Deskripsi level dari option yang dipilih (field sesuai descriptionField: descriptions/note)
  function levelDescription(row) {
    const opt = (row.levelOptions || []).find(o => o.value === row.job_management_value_id)
    return opt ? (opt[descriptionField] || '') : ''
  }

  // Cocokkan record dengan baris:
  // - matchBy 'competency' → via competency_id (baris berbasis kompetensi, mis. teknis)
  // - matchBy 'value' (default) → via job_management_value_id (baris tetap tanpa competency_id)
  function findRecord(row) {
    if (matchBy === 'competency') {
      return row.competency_id
        ? records.value.find(r => r.competency_id && r.competency_id === row.competency_id) || null
        : null
    }
    const valueIds = new Set((row.levelOptions || []).map(o => o.value))
    return records.value.find(r => r.job_management_value_id && valueIds.has(r.job_management_value_id)) || null
  }

  // Isi level terpilih dari record potency yang sudah tersimpan (via competency_id / value id per baris)
  function hydrateRows() {
    rows.value.forEach(row => {
      const rec = findRecord(row)
      row.recordId = rec ? rec.id : ''
      row.job_management_value_id = rec ? (rec.job_management_value_id || '') : ''
      // Bobot (%) — hanya terisi bila baris memakai kolom bobot (row.weight diinisialisasi)
      if (row.weight !== undefined) {
        row.weight = rec ? (rec.weight ?? row.weight) : row.weight
      }
    })
  }

  async function loadData() {
    if (!orgId.value) {
      records.value = []
      return
    }
    try {
      const res = await api.get(apiBase, { params: { organization_id: orgId.value, per_page: 100 } })
      records.value = res.data?.data || []
    } catch {
      records.value = []
    }
  }

  function askDeleteRow(row) {
    deleteTarget.value = row
    deleteError.value = ''
    deleteVisible.value = true
  }

  async function handleDelete() {
    const row = deleteTarget.value
    if (!row) return
    deleting.value = true
    deleteError.value = ''
    try {
      if (row.recordId) {
        await api.delete(`${apiBase}/${row.recordId}`)
      }
      if (afterDelete) afterDelete(row)
      deleteVisible.value = false
      await loadData()
      hydrateRows()
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 2000 })
      if (onSaved) onSaved()
    } catch (err) {
      deleteError.value = err?.response?.data?.error?.message || err.message || t('message.operation_failed')
    } finally {
      deleting.value = false
      deleteTarget.value = null
    }
  }

  async function handleSave() {
    errorMsg.value = ''
    savingCard.value = true
    try {
      for (const row of rows.value) {
        if (row.job_management_value_id) {
          // Baris tanpa competency_id → kirim hanya level (job_management_value_id)
          const payload = row.competency_id
            ? { competency_id: row.competency_id, job_management_value_id: row.job_management_value_id }
            : { job_management_value_id: row.job_management_value_id }
          // Bobot (%) opsional — dikirim hanya bila diisi
          if (row.weight !== undefined && row.weight !== null && row.weight !== '') {
            payload.weight = row.weight
          }
          if (row.recordId) {
            await api.put(`${apiBase}/${row.recordId}`, payload)
          } else {
            const res = await api.post(apiBase, { organization_id: orgId.value, ...payload })
            row.recordId = res.data?.data?.id || ''
          }
        } else if (row.recordId) {
          // Level dikosongkan → hapus record potency
          await api.delete(`${apiBase}/${row.recordId}`)
          row.recordId = ''
        }
      }
      await loadData()
      hydrateRows()
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('common.saved'), life: 2000 })
      if (onSaved) onSaved()
    } catch (err) {
      const ve = getValidationErrors(err)
      if (Object.keys(ve).length > 0) {
        errorMsg.value = Object.values(ve).join(', ')
      } else {
        errorMsg.value = err?.response?.data?.error?.message || err.message || t('message.operation_failed')
      }
    } finally {
      savingCard.value = false
    }
  }

  return {
    savingCard, errorMsg,
    deleteVisible, deleting, deleteError, deleteTarget,
    records,
    levelDescription, hydrateRows,
    loadData, askDeleteRow, handleDelete, handleSave
  }
}
