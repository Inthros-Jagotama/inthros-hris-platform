import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const apiGet = vi.fn()
const apiPost = vi.fn()
const apiPut = vi.fn()

vi.mock('@/services/api', () => ({
  default: {
    get: (...args) => apiGet(...args),
    post: (...args) => apiPost(...args),
    put: (...args) => apiPut(...args),
    delete: vi.fn(),
  },
}))

const toastAdd = vi.fn()
vi.mock('primevue/usetoast', () => ({
  useToast: () => ({ add: toastAdd }),
}))

const push = vi.fn()
const routeParams = {}
vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
  useRoute: () => ({ params: routeParams }),
}))

import DocumentTemplateForm from '@/views/settings/DocumentTemplateForm.vue'

const variablesResponse = {
  data: {
    data: [
      {
        category: 'employee',
        variables: [
          { key: 'employee.employee_id', label: 'Employee Number' },
          { key: 'employee.name', label: 'Name' },
        ],
      },
      { category: 'contract', variables: [{ key: 'contract.number', label: 'Contract Number' }] },
    ],
  },
}

// PrimeVue InputText/BaseInput membaca `$primevue.config` saat render (di app
// disuntik plugin PrimeVue di main.js) — sediakan minimal di lingkungan test.
const primevueGlobal = {
  config: {
    inputStyle: 'filled',
    inputVariant: 'filled',
    ripple: false,
    pt: {},
    ptOptions: {},
    unstyled: false,
  },
}

function mountForm() {
  return mount(DocumentTemplateForm, {
    global: {
      config: {
        globalProperties: { $primevue: primevueGlobal },
      },
    },
  })
}

function makeDocxFile(name = 'template.docx') {
  return new File(['PK\x03\x04fake-docx'], name, { type: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document' })
}

function findSaveBtn(wrapper) {
  return wrapper.findAll('button').find((b) => b.text().includes('Simpan'))
}

function setNameAndType(wrapper) {
  const nameInput = wrapper.find('input.textinput-stub')
  return nameInput.setValue('Perjanjian PKWT')
}

beforeEach(() => {
  vi.clearAllMocks()
  apiGet.mockReset()
  apiPost.mockReset()
  apiPut.mockReset()
  apiGet.mockImplementation((url) => {
    if (url.includes('/movement-types')) {
      return Promise.resolve({
        data: {
          data: [
            { value: 'promotion', label: 'Promosi' },
            { value: 'mutation', label: 'Mutasi' },
          ],
        },
      })
    }
    return Promise.resolve(variablesResponse)
  })
  delete routeParams.id
})

describe('DocumentTemplateForm (create)', () => {
  it('memuat variable reference saat mounted', async () => {
    const wrapper = mountForm()
    await flushPromises()
    expect(apiGet).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates/variables')
    expect(wrapper.text()).toContain('{{employee.name}}')
    expect(wrapper.text()).toContain('{{contract.number}}')
    // Label variabel bilingual — mengikuti bahasa aktif (id di test).
    expect(wrapper.text()).toContain('Nomor Karyawan')
    expect(wrapper.text()).toContain('Nomor Kontrak')
    expect(wrapper.text()).not.toContain('Employee Number')
  })

  it('validasi: nama kosong menampilkan error dan tidak submit', async () => {
    const wrapper = mountForm()
    await flushPromises()

    await findSaveBtn(wrapper).trigger('click')
    await flushPromises()

    expect(apiPost).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('Field ini wajib diisi')
  })

  it('validasi: tanpa file docx menampilkan error file_required', async () => {
    const wrapper = mountForm()
    await flushPromises()

    await setNameAndType(wrapper)
    const select = wrapper.find('select.select-stub')
    await select.setValue('CONTRACT_AGREEMENT')

    await findSaveBtn(wrapper).trigger('click')
    await flushPromises()

    expect(apiPost).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('File template .docx wajib diisi.')
  })

  it('tolak file non-.docx', async () => {
    const wrapper = mountForm()
    await flushPromises()

    const fileInput = wrapper.find('input[type="file"]')
    const badFile = new File(['x'], 'template.pdf', { type: 'application/pdf' })
    Object.defineProperty(fileInput.element, 'files', { value: [badFile], configurable: true })
    await fileInput.trigger('change')

    expect(wrapper.text()).toContain('Hanya file .docx yang diizinkan.')
  })

  it('tolak file lebih dari 10MB', async () => {
    const wrapper = mountForm()
    await flushPromises()

    const fileInput = wrapper.find('input[type="file"]')
    const bigFile = new File([new ArrayBuffer(11 * 1024 * 1024)], 'big.docx', {
      type: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    })
    Object.defineProperty(fileInput.element, 'files', { value: [bigFile], configurable: true })
    await fileInput.trigger('change')

    expect(wrapper.text()).toContain('10 MB')
  })

  it('submit: POST template lalu POST versi dengan FormData, lalu kembali', async () => {
    apiPost.mockResolvedValueOnce({ data: { data: { id: 'tmpl-new' } } })
    apiPost.mockResolvedValueOnce({ data: { data: { placeholders: ['employee.name'] } } })

    const wrapper = mountForm()
    await flushPromises()

    await setNameAndType(wrapper)
    const select = wrapper.find('select.select-stub')
    await select.setValue('CONTRACT_AGREEMENT')

    const fileInput = wrapper.find('input[type="file"]')
    Object.defineProperty(fileInput.element, 'files', { value: [makeDocxFile()], configurable: true })
    await fileInput.trigger('change')

    await findSaveBtn(wrapper).trigger('click')
    await flushPromises()

    expect(apiPost).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates', {
      name: 'Perjanjian PKWT',
      document_type: 'CONTRACT_AGREEMENT',
      description: undefined,
    })
    expect(apiPost).toHaveBeenCalledWith(
      '/api/v1/tenant/settings/document-templates/tmpl-new/versions',
      expect.any(FormData),
      expect.objectContaining({ headers: { 'Content-Type': 'multipart/form-data' } })
    )
    expect(toastAdd).toHaveBeenCalledWith(expect.objectContaining({ severity: 'success' }))
    expect(push).toHaveBeenCalledWith('/settings/document-templates')
  })

  it('pencarian variabel memfilter berdasarkan label (terjemahan) atau key', async () => {
    const wrapper = mountForm()
    await flushPromises()
    expect(wrapper.text()).toContain('{{employee.name}}')
    expect(wrapper.text()).toContain('{{contract.number}}')

    // Cari label Bahasa Indonesia — label terjemahan ikut dicocokkan.
    const searchInput = wrapper.find('input[placeholder*="Cari variabel"]')
    await searchInput.setValue('nomor')
    await flushPromises()

    expect(wrapper.text()).toContain('{{contract.number}}')
    expect(wrapper.text()).toContain('{{employee.employee_id}}')
    expect(wrapper.text()).not.toContain('{{employee.name}}')

    // Cari berdasarkan key — hanya variabel yang cocok yang tampil.
    await searchInput.setValue('contract')
    await flushPromises()

    expect(wrapper.text()).toContain('{{contract.number}}')
    expect(wrapper.text()).not.toContain('{{employee.name}}')

    // Tidak ada yang cocok → pesan kosong.
    await searchInput.setValue('zzzz')
    await flushPromises()
    expect(wrapper.text()).toContain('Tidak ada variabel yang cocok')

    // Bersihkan pencarian → semua variabel kembali.
    await searchInput.setValue('')
    await flushPromises()
    expect(wrapper.text()).toContain('{{employee.name}}')
  })

  it('template SK movement: pilihan jenis movement dikirim saat submit', async () => {
    apiPost.mockResolvedValueOnce({ data: { data: { id: 'tmpl-mov' } } })
    apiPost.mockResolvedValueOnce({ data: { data: { placeholders: [] } } })

    const wrapper = mountForm()
    await flushPromises()

    await setNameAndType(wrapper)
    const typeSelect = wrapper.find('select.select-stub')
    await typeSelect.setValue('MOVEMENT_SK')
    await flushPromises()

    // Dropdown jenis movement muncul setelah memilih SK Movement.
    const selects = wrapper.findAll('select.select-stub')
    expect(selects.length).toBeGreaterThanOrEqual(2)
    expect(selects[1].text()).toContain('Promosi')
    await selects[1].setValue('promotion')

    const fileInput = wrapper.find('input[type="file"]')
    Object.defineProperty(fileInput.element, 'files', { value: [makeDocxFile()], configurable: true })
    await fileInput.trigger('change')

    await findSaveBtn(wrapper).trigger('click')
    await flushPromises()

    expect(apiPost).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates', {
      name: 'Perjanjian PKWT',
      document_type: 'MOVEMENT_SK',
      movement_type: 'promotion',
      description: undefined,
    })
  })

  it('copy variable menyalin {{key}} ke clipboard + toast', async () => {
    const wrapper = mountForm()
    await flushPromises()

    const variableBtn = wrapper.findAll('button').find((b) => b.text().includes('{{employee.name}}'))
    expect(variableBtn).toBeTruthy()
    await variableBtn.trigger('click')
    await flushPromises()

    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('{{employee.name}}')
    expect(toastAdd).toHaveBeenCalledWith(expect.objectContaining({ severity: 'success' }))
  })
})

describe('DocumentTemplateForm (edit)', () => {
  it('memuat data template + file versi aktif saat mode edit', async () => {
    routeParams.id = 'tmpl-1'

    apiGet.mockImplementation((url) => {
      if (url.includes('/versions')) {
        return Promise.resolve({
          data: { data: [{ id: 'ver-1', version: 3, file_url: '/uploads/tmpl.docx', file_name: 'pkwt.docx', paper_size: 'A4', orientation: 'portrait' }] },
        })
      }
      if (url.includes('/variables')) return Promise.resolve(variablesResponse)
      return Promise.resolve({ data: { data: { id: 'tmpl-1', name: 'Perjanjian PKWT', document_type: 'CONTRACT_AGREEMENT', description: 'desc', active_version_id: 'ver-1' } } })
    })

    const wrapper = mountForm()
    await flushPromises()

    // Nama template diisi ke input (bukan text node) — cek value-nya
    expect(wrapper.find('input.textinput-stub').element.value).toBe('Perjanjian PKWT')
    expect(wrapper.text()).toContain('pkwt.docx')
  })
})
