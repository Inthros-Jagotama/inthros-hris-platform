<template>
  <div class="org-chart-container" :class="{ 'dark-mode': isDark }">
    <div ref="chartRef" class="org-chart-svg"></div>
    <div v-if="!data?.length" class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
      <i class="pi pi-sitemap text-4xl mb-2 opacity-50"></i>
      <p class="text-sm font-medium">{{ t('organization.empty_title') }}</p>
      <p class="text-sm mt-1">{{ t('organization.empty_tree') }}</p>
    </div>
  </div>
</template>
<script setup>
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from '@/composables/useI18n'
import * as d3 from 'd3'
const props = defineProps({
  data: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false }
})
const emit = defineEmits(['node-click', 'node-edit', 'node-delete', 'add-child'])
const { t } = useI18n()
const chartRef = ref(null)
let chartInstance = null
// Detect dark mode — pakai class di <html> (mengikuti app toggle)
const isDark = ref(document.documentElement.classList.contains('p-dark'))
const observer = new MutationObserver(() => {
  isDark.value = document.documentElement.classList.contains('p-dark')
})
observer.observe(document.documentElement, { attributeFilter: ['class'] })
// Flatten tree data into d3-org-chart format
function flattenForChart(nodes, parentId = null) {
  const result = []
  nodes.forEach((node, idx) => {
    result.push({
      id: node.id,
      parentId: parentId,
      code: node.code || '',
      full_code: node.full_code || '',
      nomenclature: node.nomenclature || '',
      level: node.level || 0,
      sort_order: node.sort_order || 0
    })
    if (node.children?.length) {
      result.push(...flattenForChart(node.children, node.id))
    }
  })
  return result
}
function initChart() {
  if (!chartRef.value || !props.data?.length) return
  const el = chartRef.value
  const width = el.clientWidth || 900
  // Clear previous
  el.innerHTML = ''
  const flatData = flattenForChart(props.data)
  import('d3-org-chart')
    .then(({ OrgChart }) => {
      chartInstance = new OrgChart()
      chartInstance
        .container(el)
        .data(flatData)
        .nodeWidth(() => 220)
        .nodeHeight(() => 90)
        .childrenMargin(() => 40)
        .compactMarginBetween(() => 20)
        .compactMarginPair(() => 40)
        .neighbourMargin(() => 20)
        .siblingsMargin(() => 20)
        .initialZoom(0.8)
        .nodeContent(function (d) {
          const color = d.depth === 0 ? '#059669' 
            : d.depth === 1 ? '#0ea5e9' 
            : '#6366f1'
          const textColor = '#ffffff'
          const name = String(d.data.nomenclature || '').replace(/[&<>"']/g, '')
          const code = String(d.data.code || '').replace(/[&<>"']/g, '')
          const fullCode = String(d.data.full_code || '').replace(/[&<>"']/g, '')
          return `
            <div style="background:${color};border-radius:10px;padding:12px 14px;width:210px;box-shadow:0 2px 8px rgba(0,0,0,0.12);font-family:inherit;">
              <div style="display:flex;align-items:center;gap:6px;margin-bottom:4px;">
                <span style="font-weight:600;font-size:13px;color:${textColor};white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:160px;">${name}</span>
                <span style="background:rgba(255,255,255,0.2);border-radius:4px;padding:1px 6px;font-size:10px;color:${textColor};white-space:nowrap;">${code}</span>
              </div>
              <div style="font-size:11px;color:${textColor};opacity:0.8;">
                ${fullCode}
              </div>
            </div>
          `
        })
        .onNodeClick((d) => {
          emit('node-click', d.data)
        })
        .render()
    })
    .catch((err) => {
      console.error('Failed to load OrgChart:', err)
    })
}
watch(() => props.data, () => {
  nextTick(() => {
    chartInstance = null
    initChart()
  })
})
onMounted(() => {
  nextTick(() => initChart())
})
onBeforeUnmount(() => {
  observer.disconnect()
  chartInstance = null
})
</script>
<style scoped>
.org-chart-container {
  position: relative;
  width: 100%;
  min-height: 400px;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  background: #ffffff;
  overflow: hidden;
  transition: background 0.2s, border-color 0.2s;
}
.org-chart-container.dark-mode {
  background: #1f2937;
  border-color: #374151;
}
.org-chart-svg {
  width: 100%;
  min-height: 400px;
}
.org-chart-svg :deep(.link-line) {
  stroke: #d1d5db;
  stroke-width: 2;
}
:deep(.p-datatable) {
  max-height: calc(100vh - 350px);
}
</style>
