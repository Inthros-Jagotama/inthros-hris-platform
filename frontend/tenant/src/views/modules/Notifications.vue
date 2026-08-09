<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
      <div class="flex items-center gap-2">
        <Button :label="t('notification.filter_all')" size="small" :severity="filterUnreadOnly ? 'secondary' : 'primary'" :outlined="filterUnreadOnly" @click="setFilter(false)" />
        <Button :label="t('notification.filter_unread')" size="small" :severity="filterUnreadOnly ? 'primary' : 'secondary'" :outlined="!filterUnreadOnly" @click="setFilter(true)" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <Button :label="t('notification.mark_all_read')" icon="pi pi-check" size="small" severity="secondary" outlined :disabled="notifState.unreadCount === 0" @click="handleMarkAllAsRead" />
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
    <DataTable
      v-else
      :value="items"
      lazy
      :totalRecords="totalRecords"
      :first="firstRecord"
      :rows="perPage"
      @page="onPage($event)"
      paginator
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 20, 50]"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
      rowHover
      @row-click="handleRowClick($event.data)"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-bell-slash text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('notification.empty') }}</p>
        </div>
      </template>
      <Column style="width:24px">
        <template #body="{data}">
          <span v-if="!data.is_read" class="block w-2 h-2 rounded-full bg-emerald-500"></span>
        </template>
      </Column>
      <Column field="title" :header="t('notification.column_title')">
        <template #body="{data}">
          <span class="font-medium text-gray-800 dark:text-gray-100" :class="{ 'font-semibold': !data.is_read }">{{ data.title }}</span>
        </template>
      </Column>
      <Column field="body" :header="t('notification.column_body')">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ data.body }}</span>
        </template>
      </Column>
      <Column field="type" :header="t('notification.column_type')" style="width:180px">
        <template #body="{data}">
          <Tag :value="data.type" severity="secondary" class="!text-xs" />
        </template>
      </Column>
      <Column field="created_at" :header="t('notification.column_date')" style="width:170px">
        <template #body="{data}">
          <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatDate(data.created_at) }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <Button
            v-if="!data.is_read"
            icon="pi pi-check"
            size="small"
            text
            v-tooltip.left="t('notification.mark_read')"
            @click.stop="handleMarkAsRead(data)"
          />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { formatDate as formatDateGlobal } from '@/utils/formatDate'
import { useNotifications } from '@/stores/notifications'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import SkeletonTable from '@/components/SkeletonTable.vue'

const { t, locale } = useI18n()
const router = useRouter()
const { state: notifState, markAsRead, markAllAsRead, refresh: refreshBadge } = useNotifications()

function formatDate(v) {
  if (!v) return '-'
  const datePart = formatDateGlobal(v, locale.value)
  if (!datePart) return '-'
  const time = new Date(v).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  return `${datePart} ${time}`
}

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(20)
const filterUnreadOnly = ref(false)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const skeletonColumns = [
  { type: 'text', width: 'w-4', headerWidth: 'w-4' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-56', headerWidth: 'w-32' },
  { type: 'tag', width: 'w-28', headerWidth: 'w-20' },
  { type: 'text', width: 'w-28', headerWidth: 'w-24' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

// NOTE: this endpoint's envelope nests the paginated payload one level
// deeper than most other list endpoints: { success, data: { data, total,
// page, per_page } } instead of { success, data, total, page, per_page }.
async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (filterUnreadOnly.value) params.is_read = false
    const res = await api.get('/api/v1/tenant/notifications', { params })
    const body = res.data?.data || {}
    items.value = body.data || []
    totalRecords.value = body.total || 0
    if (body.page) currentPage.value = body.page
  } catch {
    items.value = []
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function setFilter(unreadOnly) {
  filterUnreadOnly.value = unreadOnly
  currentPage.value = 1
  loadData()
}

async function handleMarkAsRead(item) {
  await markAsRead(item.id)
  await loadData()
}

async function handleMarkAllAsRead() {
  await markAllAsRead()
  await loadData()
}

// Clicking a row marks it read and (when it points at a module the FE
// already has a page for) navigates there. Reference types without a
// concrete FE destination yet are intentionally left as "mark read only" —
// no speculative deep-link map for modules that are still placeholders.
async function handleRowClick(item) {
  if (!item.is_read) {
    await markAsRead(item.id)
    await loadData()
  }
  if (item.reference_type === 'leave') {
    router.push('/leave')
    return
  }
  // Template KPI/OKR baru dibuat → langsung ke halaman "isi KPI/OKR saya"
  // (self-assessment) agar karyawan langsung mengisi evaluasinya.
  if (item.type === 'KPI_TEMPLATE_CREATED' || item.reference_type === 'performance_kpi_template') {
    router.push('/performance/kpi/my-evaluation')
    return
  }
  if (item.type === 'OKR_TEMPLATE_CREATED' || item.reference_type === 'okr_template') {
    router.push('/performance/okr/my-evaluation')
    return
  }
  // Lembur (dua alur §32b): notifikasi OVERTIME_ASSIGNED (penugasan),
  // OVERTIME_ACTUAL_APPROVED/REJECTED (hasil approval aktual) → halaman lembur.
  if (item.reference_type === 'attendance_overtime' ||
      item.type?.startsWith('OVERTIME_') ||
      item.type?.startsWith('OVERTIME_ACTUAL_')) {
    router.push('/attendance/overtime')
    return
  }
  // Career movement (modul employeemovement): MOVEMENT_SUBMITTED/APPROVED/
  // REJECTED/EXECUTED → halaman movement.
  if (item.reference_type === 'employeemovement' || item.type?.startsWith('MOVEMENT_')) {
    router.push('/admin/career/movements')
    return
  }
}

onMounted(() => {
  loadData()
  refreshBadge()
})
</script>
