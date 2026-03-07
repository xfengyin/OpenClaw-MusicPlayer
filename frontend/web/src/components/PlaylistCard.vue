<template>
  <div class="playlist-card" @click="handleClick">
    <div class="cover-wrapper">
      <img :src="cover" :alt="title" class="cover" />
      <div class="play-overlay">
        <n-icon size="32" color="#fff"><play-circle /></n-icon>
      </div>
      <div class="play-count">
        <n-icon size="12"><headset /></n-icon>
        <span>{{ formatPlayCount }}</span>
      </div>
    </div>
    <div class="title">{{ title }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NIcon } from 'naive-ui'
import { PlayCircle, Headset } from '@vicons/ionicons5'

const props = defineProps<{
  title: string
  cover: string
  playCount: number
}>()

const emit = defineEmits<{
  click: []
}>()

const formatPlayCount = computed(() => {
  if (props.playCount >= 10000) {
    return (props.playCount / 10000).toFixed(1) + '万'
  }
  return props.playCount.toString()
})

const handleClick = () => {
  emit('click')
}
</script>

<style scoped>
.playlist-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.playlist-card:hover {
  transform: translateY(-4px);
}

.cover-wrapper {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  aspect-ratio: 1;
}

.cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.play-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.playlist-card:hover .play-overlay {
  opacity: 1;
}

.play-count {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 12px;
  font-size: 12px;
  color: #fff;
}

.title {
  margin-top: 8px;
  font-size: 14px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
