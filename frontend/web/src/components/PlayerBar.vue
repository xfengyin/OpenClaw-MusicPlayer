<template>
  <div class="player-bar">
    <!-- Song Info -->
    <div class="song-info">
      <n-avatar
        :size="56"
        :src="currentSong.cover || '/default-cover.png'"
        shape="square"
      />
      <div class="song-details">
        <div class="song-title">{{ currentSong.title || '未播放' }}</div>
        <div class="song-artist">{{ currentSong.artist || '-' }}</div>
      </div>
    </div>

    <!-- Controls -->
    <div class="controls">
      <div class="control-buttons">
        <n-button quaternary circle size="small">
          <n-icon size="20"><play-skip-back /></n-icon>
        </n-button>
        <n-button circle type="primary" size="large" @click="togglePlay">
          <n-icon size="24">
            <play v-if="!isPlaying" />
            <pause v-else />
          </n-icon>
        </n-button>
        <n-button quaternary circle size="small">
          <n-icon size="20"><play-skip-forward /></n-icon>
        </n-button>
      </div>
      <div class="progress-bar">
        <span class="time">{{ formatTime(currentTime) }}</span>
        <n-slider
          v-model:value="progress"
          :max="100"
          :step="0.1"
          class="slider"
        />
        <span class="time">{{ formatTime(duration) }}</span>
      </div>
    </div>

    <!-- Volume & Settings -->
    <div class="extra-controls">
      <n-button quaternary circle size="small">
        <n-icon size="18"><repeat /></n-icon>
      </n-button>
      <n-button quaternary circle size="small">
        <n-icon size="18"><list /></n-icon>
      </n-button>
      <div class="volume-control">
        <n-icon size="18"><volume-high /></n-icon>
        <n-slider v-model:value="volume" :max="100" class="volume-slider" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  NAvatar,
  NButton,
  NIcon,
  NSlider
} from 'naive-ui'
import {
  PlaySkipBack,
  PlaySkipForward,
  Play,
  Pause,
  Repeat,
  List,
  VolumeHigh
} from '@vicons/ionicons5'

const currentSong = ref({
  title: '',
  artist: '',
  cover: ''
})

const isPlaying = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const progress = ref(0)
const volume = ref(80)

const togglePlay = () => {
  isPlaying.value = !isPlaying.value
}

const formatTime = (seconds: number): string => {
  if (!seconds) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}
</script>

<style scoped>
.player-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 24px;
  background: rgba(255, 255, 255, 0.05);
}

.song-info {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 200px;
}

.song-details {
  overflow: hidden;
}

.song-title {
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.song-artist {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.controls {
  flex: 1;
  max-width: 500px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.control-buttons {
  display: flex;
  align-items: center;
  gap: 16px;
}

.progress-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.time {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  min-width: 40px;
  text-align: center;
}

.slider {
  flex: 1;
}

.extra-controls {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 200px;
  justify-content: flex-end;
}

.volume-control {
  display: flex;
  align-items: center;
  gap: 8px;
}

.volume-slider {
  width: 80px;
}
</style>
