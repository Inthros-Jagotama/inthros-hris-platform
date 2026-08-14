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
      { category: 'employee', variables: [{ key: 'employee.name', label: 'Name' }] },
      { category: 'contract', variables: [{ key: 'contract.number', label: 'Contract Number' }] },
    ],
  },
}

function mountForm() {
  return mount(DocumentTemplateForm)
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
  apiGet.mockResolvedValue(variablesResponse)
  delete routeParams.id
})

describe('DocumentTemplateForm (create)', () => {
  it('memuat variable reference saat mounted', async () => {
    const wrapper = mountForm()
    await flushPromises()
    expect(apiGet).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates/variables')
    expect(wrapper.text()).toContain('{{employee.name}}')
    expect(wrapper.text()).toContain('{{contract.number}}')
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
