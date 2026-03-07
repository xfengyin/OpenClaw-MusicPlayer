<template>
  <div class="player-container">
    <!-- Audio Element -->
    <audio
      ref="audioRef"
      :src="currentSong?.url"
      @timeupdate="onTimeUpdate"
      @ended="onEnded"
      @loadedmetadata="onLoadedMetadata"
      @error="onError"
    />

    <!-- Player Controls -->
    <PlayerBar
      :current-song="currentSong"
      :is-playing="isPlaying"
      :current-time="currentTime"
      :duration="duration"
      :volume="volume"
      @toggle-play="togglePlay"
      @seek="seek"
      @volume-change="setVolume"
      @next="next"
      @prev="prev"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { usePlayerStore } from '../stores/player'
import PlayerBar from './PlayerBar.vue'

const playerStore = usePlayerStore()
const audioRef = ref<HTMLAudioElement>()

const currentSong = computed(() => playerStore.currentSong)
const isPlaying = computed(() => playerStore.isPlaying)
const currentTime = computed(() => playerStore.currentTime)
const duration = computed(() => playerStore.currentSong?.duration || 0)
const volume = computed(() => playerStore.volume)

// Watch for play state changes
watch(isPlaying, (playing) => {
  if (!audioRef.value) return
  
  if (playing) {
    audioRef.value.play().catch(err => {
      console.error('Play error:', err)
      playerStore.pause()
    })
  } else {
    audioRef.value.pause()
  }
})

// Watch for song changes
watch(currentSong, (song) => {
  if (!song || !audioRef.value) return
  
  // Reset time when song changes
  playerStore.setCurrentTime(0)
  
  // Auto play if enabled
  if (isPlaying.value) {
    audioRef.value.play().catch(err => {
      console.error('Auto play error:', err)
      playerStore.pause()
    })
  }
})

// Event handlers
const onTimeUpdate = () => {
  if (!audioRef.value) return
  playerStore.setCurrentTime(audioRef.value.currentTime)
}

const onEnded = () => {
  playerStore.next()
}

const onLoadedMetadata = () => {
  if (!audioRef.value) return
  // Update duration if different
  if (audioRef.value.duration !== duration.value) {
    // Could update store here if needed
  }
}

const onError = (e: Event) => {
  console.error('Audio error:', e)
  playerStore.pause()
}

// Control methods
const togglePlay = () => {
  playerStore.togglePlay()
}

const seek = (time: number) => {
  if (!audioRef.value) return
  audioRef.value.currentTime = time
  playerStore.setCurrentTime(time)
}

const setVolume = (vol: number) => {
  if (!audioRef.value) return
  audioRef.value.volume = vol / 100
  playerStore.setVolume(vol)
}

const next = () => {
  playerStore.next()
}

const prev = () => {
  playerStore.prev()
}
</script>

<style scoped>
.player-container {
  width: 100%;
}
</style>
