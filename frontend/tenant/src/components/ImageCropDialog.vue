<template>
  <Dialog
    v-model:visible="visibleModel"
    :header="title"
    modal
    :style="{ width: '560px' }"
    :closable="!cropping"
    :dismissable-mask="false"
    :close-on-escape="!cropping"
    @hide="onHide"
  >
    <div class="space-y-3 flex flex-col" style="height:70vh">
      <!-- Crop area → fills remaining space after aspect-ratio + zoom controls -->
      <div class="relative bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden flex-1" style="min-height:200px">
        <img ref="imageRef" :src="imageSrc" alt="Crop" class="max-w-full" />
      </div>

      <!-- Aspect ratio options -->
      <div class="flex items-center gap-2 flex-wrap">
        <span class="text-xs font-medium text-gray-500 dark:text-gray-400 mr-1">{{ t('common.aspect_ratio') }}:</span>
        <Button v-for="opt in aspectOptions" :key="opt.key"
          :label="opt.label"
          size="small"
          :severity="selectedAspect === opt.key ? 'primary' : 'secondary'"
          :outlined="selectedAspect !== opt.key"
          @click="setAspect(opt.key)" />
        <span v-if="customAspect" class="flex items-center gap-1 ml-1">
          <InputText v-model="aspectW" size="small" class="!w-12 !text-xs !py-1 !px-1.5 text-center" placeholder="W" />
          <span class="text-xs text-gray-400">:</span>
          <InputText v-model="aspectH" size="small" class="!w-12 !text-xs !py-1 !px-1.5 text-center" placeholder="H" />
          <Button icon="pi pi-check" size="small" text severity="primary" @click="applyCustomAspect" />
        </span>
      </div>

      <!-- Zoom control -->
      <div class="flex items-center gap-2">
        <i class="pi pi-search-minus text-gray-400 text-sm"></i>
        <input type="range" min="0" max="100" :value="zoomValue" @input="onZoom" class="flex-1 accent-emerald-500 h-1 cursor-pointer" />
        <i class="pi pi-search-plus text-gray-400 text-sm"></i>
        <Button icon="pi pi-undo" size="small" text severity="secondary" v-tooltip.top="t('common.reset')" @click="resetCrop" />
      </div>
    </div>

    <template #footer>
      <div class="flex items-center justify-end gap-2">
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" :disabled="cropping" @click="onCancel" />
        <Button :label="t('common.apply')" icon="pi pi-check" size="small" :loading="cropping" :disabled="cropping" @click="onCrop" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick, onBeforeUnmount } from 'vue'
import { useI18n } from '@/composables/useI18n'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Cropper from 'cropperjs'

const { t } = useI18n()

const props = defineProps({
  visible: { type: Boolean, default: false },
  imageSrc: { type: String, default: '' },
  title: { type: String, default: 'Crop Image' }
})

const emit = defineEmits(['update:visible', 'crop', 'cancel'])

const visibleModel = computed({
  get: () => props.visible,
  set: (val) => {
    if (!cropping.value) {
      emit('update:visible', val)
    }
  }
})

const imageRef = ref(null)
let cropper = null
const cropping = ref(false)
const zoomValue = ref(50)
let prevZoomVal = 50 // for relative zoom delta tracking

// Aspect ratio options
const selectedAspect = ref('free')
const customAspect = ref(false)
const aspectW = ref('')
const aspectH = ref('')

const aspectOptions = [
  { key: 'free', label: 'Free' },
  { key: '1:1', label: '1:1' },
  { key: '4:3', label: '4:3' },
  { key: '16:9', label: '16:9' },
  { key: '3:2', label: '3:2' },
  { key: 'custom', label: t('common.custom') }
]

function getAspectRatio(key) {
  switch (key) {
    case '1:1': return 1
    case '4:3': return 4 / 3
    case '16:9': return 16 / 9
    case '3:2': return 3 / 2
    default: return NaN
  }
}

function setAspectRatio(ratio) {
  const selection = cropper.getCropperSelection()
  if (!selection) return
  if (isNaN(ratio)) {
    selection.removeAttribute('aspect-ratio')
  } else {
    selection.setAttribute('aspect-ratio', String(ratio))
  }
}

function setAspect(key) {
  if (key === 'custom') {
    customAspect.value = true
    selectedAspect.value = 'custom'
    return
  }
  customAspect.value = false
  selectedAspect.value = key
  if (cropper) {
    setAspectRatio(getAspectRatio(key))
  }
}

function applyCustomAspect() {
  const w = parseFloat(aspectW.value)
  const h = parseFloat(aspectH.value)
  if (w > 0 && h > 0 && cropper) {
    setAspectRatio(w / h)
  }
}

function initCropper() {
  destroyCropper()
  if (!imageRef.value) return
  nextTick(() => {
    cropper = new Cropper(imageRef.value, {
      viewMode: 1,
      dragMode: 'move',
      initialCenterSize: 'cover',
      initialAspectRatio: NaN,
      aspectRatio: NaN,
      minCropBoxWidth: 100,
      minCropBoxHeight: 100,
      zoomOnWheel: true,
      zoomable: true,
      rotatable: true,
      scalable: false,
      cropBoxMovable: true,
      cropBoxResizable: true,
      toggleDragModeOnDblclick: false
    })
    // Force cropper-canvas to fill parent div height
    const canvasEl = cropper.getCropperCanvas()
    if (canvasEl) {
      canvasEl.style.height = '100%'
    }
    zoomValue.value = 50
    prevZoomVal = 50
  })
}

function destroyCropper() {
  if (cropper) {
    cropper.destroy()
    cropper = null
  }
}function onZoom(e) {
  if (!cropper) return
  const val = parseFloat(e.target.value)
  zoomValue.value = val

  // cropperjs v2 $zoom() is RELATIVE — compute delta from previous slider position
  const prevRatio =
    prevZoomVal <= 50
      ? 0.1 + (prevZoomVal / 50) * 0.9
      : 1.0 + ((prevZoomVal - 50) / 50) * 2.0
  const newRatio =
    val <= 50 ? 0.1 + (val / 50) * 0.9 : 1.0 + ((val - 50) / 50) * 2.0

  const factor = newRatio / prevRatio
  // Convert multiplicative factor to $zoom() scale parameter
  const scale = factor >= 1 ? factor - 1 : 1 - 1 / factor

  const image = cropper.getCropperImage()
  if (image) image.$zoom(scale)

  prevZoomVal = val
}

function resetCrop() {
  if (!cropper) return
  const image = cropper.getCropperImage()
  if (image) {
    image.$resetTransform()
    image.$center()
  }
  zoomValue.value = 50
  prevZoomVal = 50
}

async function onCrop() {
  if (!cropper) return
  cropping.value = true
  try {
    const selection = cropper.getCropperSelection()
    if (!selection) { cropping.value = false; return }

    const canvas = await selection.$toCanvas({
      width: 512,
      height: 512
    })
    canvas.toBlob((blob) => {
      emit('crop', blob)
      cropping.value = false
      emit('update:visible', false)
    }, 'image/jpeg', 0.92)
  } catch (e) {
    cropping.value = false
    console.error('Crop failed:', e)
  }
}

function onCancel() {
  emit('cancel')
  emit('update:visible', false)
}

function onHide() {
  if (cropping.value) return
  destroyCropper()
  // Emit cancel so parent can clear image src
  emit('cancel')
}

watch(() => props.visible, (val) => {
  if (val) {
    nextTick(() => initCropper())
  } else {
    destroyCropper()
  }
})

onBeforeUnmount(() => {
  destroyCropper()
})
</script>
