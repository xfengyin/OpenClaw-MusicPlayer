<template>
  <div class="settings">
    <n-h1>设置</n-h1>
    
    <n-card title="下载设置">
      <n-form label-placement="left" label-width="auto">
        <n-form-item label="下载路径">
          <n-input v-model:value="downloadPath" placeholder="选择下载路径" />
        </n-form-item>
        <n-form-item label="音质选择">
          <n-select v-model:value="quality" :options="qualityOptions" />
        </n-form-item>
        <n-form-item label="同时下载数">
          <n-slider v-model:value="maxDownloads" :max="10" />
        </n-form-item>
      </n-form>
    </n-card>

    <n-card title="播放设置" style="margin-top: 16px;">
      <n-form label-placement="left" label-width="auto">
        <n-form-item label="自动播放">
          <n-switch v-model:value="autoPlay" />
        </n-form-item>
        <n-form-item label="淡入淡出">
          <n-switch v-model:value="fadeEffect" />
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  NH1,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NSwitch,
  NSlider
} from 'naive-ui'

const downloadPath = ref('')
const quality = ref('high')
const maxDownloads = ref(3)
const autoPlay = ref(true)
const fadeEffect = ref(true)

const qualityOptions = [
  { label: '标准 (128kbps)', value: 'standard' },
  { label: '较高 (192kbps)', value: 'medium' },
  { label: '极高 (320kbps)', value: 'high' },
  { label: '无损 (FLAC)', value: 'lossless' }
]
</script>

<style scoped>
.settings {
  max-width: 800px;
}

.n-h1 {
  margin-bottom: 24px;
}
</style>
