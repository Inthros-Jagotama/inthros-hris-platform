import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const apiGet = vi.fn()

vi.mock('@/services/api', () => ({
  default: {
    get: (...args) => apiGet(...args),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

const toastAdd = vi.fn()
vi.mock('primevue/usetoast', () => ({
  useToast: () => ({ add: toastAdd }),
}))

const push = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push, replace: push }),
  useRoute: () => ({ params: { id: 'emp-1' } }),
}))

vi.mock('@/stores/auth', () => ({
  useAuth: () => ({ hasPermission: () => true }),
}))

import EmployeeDetail from '@/views/modules/employee/EmployeeDetail.vue'

const employee = {
  id: 'emp-1',
  employee_id: 'EMP-001',
  name: 'Budi Santoso',
  nik: '3201010101',
  gender: 'M',
  status: 'active',
  religion_id: 'rel-1',
  dob: '1995-06-15',
  addresses: [{ type: 'MAIN', address: 'Jl. Merdeka 1', village_id: 'v1', postal_code: '10110' }],
  emergency_contacts: [{ name: 'Siti', phone_number: '0812', relationship_type_id: 'rt-1' }],
  families: [{ name: 'Ani', relationship_type_id: 'rt-1', dob: '1998-02-20' }],
  educations: [],
  experiences: [],
  documents: [{ name: 'KTP.pdf', file: '/uploads/documents/emp-1_abc.pdf', note: 'Foto KTP' }],
  insurances: [],
  banks: [],
  employments: [],
}

const emptyList = { data: { data: [], total: 0, page: 1 } }

const careerHistory = {
  data: {
    data: {
      employee_id: 'emp-1',
      employee_name: 'Budi Santoso',
      employee_code: 'EMP-001',
      current_position: {
        employment_id: 'emp1',
        effective_date: '2025-03-01',
        position_name: 'Senior Software Engineer',
        organization_name: 'Engineering',
        employment_status_name: 'Permanent',
      },
      timeline: [
        { date: '2024-01-01', event_type: 'JOINED', title: 'Software Engineer', description: 'Bergabung sebagai Software Engineer' },
        { date: '2025-03-01', event_type: 'MOVEMENT', title: '', description: 'Software Engineer → Senior Software Engineer', movement_type: 'promotion' },
        { date: '2025-03-01', event_type: 'CONTRACT', title: 'KTR-001', contract_type: 'pkwt' },
      ],
    },
  },
}

function mockApiByUrl() {
  apiGet.mockImplementation((url) => {
    if (url.includes('career-history')) return Promise.resolve(careerHistory)
    if (url.includes('/employees/')) return Promise.resolve({ data: { data: employee } })
    if (url.includes('religions')) return Promise.resolve({ data: { data: [{ id: 'rel-1', name: 'Islam' }] } })
    if (url.includes('marital-statuses')) return Promise.resolve({ data: { data: [] } })
    if (url.includes('nationalities')) return Promise.resolve({ data: { data: [] } })
    if (url.includes('relationship-types')) return Promise.resolve({ data: { data: [{ id: 'rt-1', name: 'Saudara' }] } })
    if (url.includes('educations')) return Promise.resolve({ data: { data: [] } })
    if (url.includes('insurances')) return Promise.resolve({ data: { data: [] } })
    if (url.includes('banks')) return Promise.resolve({ data: { data: [] } })
    if (url.includes('organizations')) return Promise.resolve({ data: { data: [] } })
    if (url.includes('employment-statuses')) return Promise.resolve({ data: { data: [] } })
    if (url.includes('villages/')) return Promise.resolve({ data: { data: { name: 'Kebon Kelapa', district_name: 'Gambir', regency_name: 'Jakarta Pusat', province_name: 'DKI Jakarta' } } })
    if (url.includes('employee-payroll-profiles')) return Promise.resolve({ data: { data: [], total: 0 } })
    if (url.includes('employee-bank-profiles')) return Promise.resolve({ data: { data: [], total: 0 } })
    if (url.includes('employee-bpjs-profiles')) return Promise.resolve({ data: { data: [], total: 0 } })
    if (url.includes('employee-tax-profiles')) return Promise.resolve({ data: { data: [], total: 0 } })
    return Promise.resolve({ data: { data: [] } })
  })
}

function findSubmenu(wrapper, text) {
  return wrapper.findAll('[role="button"]').find((w) => w.text().includes(text))
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('EmployeeDetail', () => {
  it('merender header, sub-menu, dan panel personal setelah load', async () => {
    mockApiByUrl()
    const wrapper = mount(EmployeeDetail)
    await flushPromises()

    expect(wrapper.text()).toContain('Budi Santoso')
    expect(wrapper.text()).toContain('Profil')
    // sub-menu berisi semua item
    expect(findSubmenu(wrapper, 'Alamat')).toBeTruthy()
    expect(findSubmenu(wrapper, 'Profil Payroll')).toBeTruthy()
    // nilai dari label ref data muncul (agama)
    expect(wrapper.text()).toContain('Islam')
    // tanggal lahir diformat (bukan mentah YYYY-MM-DD)
    expect(wrapper.text()).toContain('15 Juni 1995')
    expect(wrapper.text()).not.toContain('1995-06-15')
  })

  it('klik sub-menu menampilkan panel section terkait', async () => {
    mockApiByUrl()
    const wrapper = mount(EmployeeDetail)
    await flushPromises()

    const alamat = findSubmenu(wrapper, 'Alamat')
    await alamat.trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Jl. Merdeka 1')
    expect(wrapper.text()).toContain('Kebon Kelapa')

    const kontak = findSubmenu(wrapper, 'Kontak')
    await kontak.trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('Siti')
  })

  it('panel Dokumen menyediakan link download file', async () => {
    mockApiByUrl()
    const wrapper = mount(EmployeeDetail)
    await flushPromises()

    await findSubmenu(wrapper, 'Dokumen').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('KTP.pdf')
    const link = wrapper.find('a[href="/uploads/documents/emp-1_abc.pdf"]')
    expect(link.exists()).toBe(true)
    expect(link.attributes('download')).toBe('KTP.pdf')
  })

  it('panel Career Timeline menampilkan posisi saat ini + event timeline', async () => {
    mockApiByUrl()
    const wrapper = mount(EmployeeDetail)
    await flushPromises()

    await findSubmenu(wrapper, 'Linimasa Karier').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Posisi Saat Ini')
    expect(wrapper.text()).toContain('Senior Software Engineer')
    expect(wrapper.text()).toContain('Bergabung')
    expect(wrapper.text()).toContain('Promosi')
    expect(wrapper.text()).toContain('Kontrak')
    expect(wrapper.text()).toContain('Software Engineer → Senior Software Engineer')
  })

  it('tetap tampil walau endpoint payroll gagal (tidak redirect)', async () => {
    apiGet.mockImplementation((url) => {
      if (url.includes('/employees/')) return Promise.resolve({ data: { data: employee } })
      if (url.includes('employee-payroll-profiles') || url.includes('employee-bank-profiles')
        || url.includes('employee-bpjs-profiles') || url.includes('employee-tax-profiles')) {
        return Promise.reject(new Error('forbidden'))
      }
      return Promise.resolve({ data: { data: [] } })
    })
    const wrapper = mount(EmployeeDetail)
    await flushPromises()

    expect(wrapper.text()).toContain('Budi Santoso')
    expect(push).not.toHaveBeenCalled()
  })
})
