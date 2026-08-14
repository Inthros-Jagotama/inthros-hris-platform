import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const apiGet = vi.fn()
const apiPost = vi.fn()
const apiDelete = vi.fn()
const apiPut = vi.fn()

vi.mock('@/services/api', () => ({
  default: {
    get: (...args) => apiGet(...args),
    post: (...args) => apiPost(...args),
    put: (...args) => apiPut(...args),
    delete: (...args) => apiDelete(...args),
  },
}))

const toastAdd = vi.fn()
vi.mock('primevue/usetoast', () => ({
  useToast: () => ({ add: toastAdd }),
}))

const push = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
  useRoute: () => ({ params: {} }),
}))

import DocumentTemplatesView from '@/views/settings/DocumentTemplatesView.vue'

const sampleTemplate = {
  id: 'tmpl-1',
  name: 'Perjanjian PKWT',
  code: 'TMPL-CONTRACT-ABC123',
  document_type: 'CONTRACT_AGREEMENT',
  status: 'ACTIVE',
  active_version_id: 'ver-1',
  updated_at: '2026-08-01T10:00:00Z',
}

const listResponse = {
  data: {
    data: {
      data: [sampleTemplate],
      total: 1,
      page: 1,
    },
  },
}

const emptyResponse = { data: { data: { data: [], total: 0, page: 1 } } }

function mountView() {
  return mount(DocumentTemplatesView)
}

// Cari button aksi berdasarkan icon class di dalamnya
function findActionBtn(wrapper, iconClass) {
  return wrapper.findAll('button').find((b) => b.find(`i.${iconClass}`).exists())
}

beforeEach(() => {
  vi.clearAllMocks()
  apiGet.mockReset()
  apiPost.mockReset()
  apiDelete.mockReset()
  apiPut.mockReset()
})

describe('DocumentTemplatesView', () => {
  it('memuat daftar template saat mounted', async () => {
    apiGet.mockResolvedValue(listResponse)
    const wrapper = mountView()
    await flushPromises()

    expect(apiGet).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates', expect.objectContaining({ params: expect.any(Object) }))
    expect(wrapper.text()).toContain('Perjanjian PKWT')
    expect(wrapper.text()).toContain('TMPL-CONTRACT-ABC123')
  })

  it('menampilkan pesan kosong bila tidak ada data', async () => {
    apiGet.mockResolvedValue(emptyResponse)
    const wrapper = mountView()
    await flushPromises()
    expect(wrapper.text()).toContain('Belum ada template dokumen')
  })

  it('openDetail memuat detail + versions lalu menampilkan dialog', async () => {
    apiGet.mockResolvedValue(listResponse)
    const wrapper = mountView()
    await flushPromises()

    apiGet.mockImplementation((url) => {
      if (url.includes('/versions')) {
        return Promise.resolve({ data: { data: [{ id: 'ver-1', version: 3, content: '<p>x</p>', paper_size: 'A4', orientation: 'portrait' }] } })
      }
      return Promise.resolve({ data: { data: sampleTemplate } })
    })

    const eyeBtn = findActionBtn(wrapper, 'pi-eye')
    expect(eyeBtn).toBeTruthy()
    await eyeBtn.trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Perjanjian PKWT')
    expect(wrapper.text()).toContain('TMPL-CONTRACT-ABC123')
  })

  it('goCreate menavigasi ke halaman form baru', async () => {
    apiGet.mockResolvedValue(listResponse)
    const wrapper = mountView()
    await flushPromises()

    const newBtn = wrapper.findAll('button').find((b) => b.text().includes('Template Baru'))
    await newBtn.trigger('click')
    expect(push).toHaveBeenCalledWith('/settings/document-templates/new')
  })

  it('preview memanggil endpoint preview dan menampilkan PDF', async () => {
    apiGet.mockResolvedValue(listResponse)
    const wrapper = mountView()
    await flushPromises()

    apiPost.mockResolvedValue({
      data: { data: { pdf_url: '/uploads/previews/test.pdf', file_name: 'preview_test.pdf' } },
    })

    const pdfBtn = findActionBtn(wrapper, 'pi-file-pdf')
    await pdfBtn.trigger('click')
    await flushPromises()

    expect(apiPost).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates/tmpl-1/preview')
    const iframe = wrapper.find('iframe')
    expect(iframe.exists()).toBe(true)
    expect(iframe.attributes('src')).toBe('/uploads/previews/test.pdf')
  })

  it('preview menampilkan error bila gagal', async () => {
    apiGet.mockResolvedValue(listResponse)
    const wrapper = mountView()
    await flushPromises()

    apiPost.mockRejectedValue({ response: { data: { error: { message: 'gagal' } } } })

    const pdfBtn = findActionBtn(wrapper, 'pi-file-pdf')
    await pdfBtn.trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('gagal')
  })

  it('delete: konfirmasi lalu DELETE + reload list', async () => {
    apiGet.mockResolvedValue(listResponse)
    const wrapper = mountView()
    await flushPromises()

    apiDelete.mockResolvedValue({ data: {} })
    // reload setelah delete
    apiGet.mockResolvedValue(emptyResponse)

    const trashBtn = findActionBtn(wrapper, 'pi-trash')
    await trashBtn.trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Hapus template')
    expect(wrapper.text()).toContain('Perjanjian PKWT')
    await wrapper.find('.confirm-ok').trigger('click')
    await flushPromises()

    expect(apiDelete).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates/tmpl-1')
    expect(toastAdd).toHaveBeenCalledWith(expect.objectContaining({ severity: 'success' }))
  })

  it('activate: konfirmasi lalu POST activate', async () => {
    const inactive = { ...sampleTemplate, id: 'tmpl-2', status: 'INACTIVE' }
    apiGet.mockResolvedValue({ data: { data: { data: [inactive], total: 1, page: 1 } } })
    const wrapper = mountView()
    await flushPromises()

    apiPost.mockResolvedValue({ data: {} })
    apiGet.mockResolvedValue(listResponse)

    const activateBtn = findActionBtn(wrapper, 'pi-check-circle')
    expect(activateBtn).toBeTruthy()
    await activateBtn.trigger('click')
    await flushPromises()

    await wrapper.find('.confirm-ok').trigger('click')
    await flushPromises()

    expect(apiPost).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates/tmpl-2/activate')
  })

  it('deactivate: konfirmasi lalu POST deactivate', async () => {
    apiGet.mockResolvedValue(listResponse)
    const wrapper = mountView()
    await flushPromises()

    apiPost.mockResolvedValue({ data: {} })

    const deactivateBtn = findActionBtn(wrapper, 'pi-times-circle')
    expect(deactivateBtn).toBeTruthy()
    await deactivateBtn.trigger('click')
    await flushPromises()

    await wrapper.find('.confirm-ok').trigger('click')
    await flushPromises()

    expect(apiPost).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates/tmpl-1/deactivate')
  })

  it('openVersions menampilkan dialog daftar versi', async () => {
    apiGet.mockResolvedValue(listResponse)
    const wrapper = mountView()
    await flushPromises()

    apiGet.mockImplementation((url) => {
      if (url.includes('/versions')) {
        return Promise.resolve({
          data: { data: [{ id: 'ver-1', version: 3, paper_size: 'A4', orientation: 'portrait', created_at: '2026-08-01T10:00:00Z' }] },
        })
      }
      return Promise.resolve(listResponse)
    })

    const versionBtn = findActionBtn(wrapper, 'pi-history')
    await versionBtn.trigger('click')
    await flushPromises()

    expect(apiGet).toHaveBeenCalledWith('/api/v1/tenant/settings/document-templates/tmpl-1/versions')
    expect(wrapper.text()).toContain('v3')
  })
})
