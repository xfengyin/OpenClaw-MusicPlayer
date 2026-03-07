import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Song {
  id: string
  title: string
  artist: string
  album: string
  cover?: string
  duration: number
  url?: string
}

export const usePlayerStore = defineStore('player', () => {
  // State
  const currentSong = ref<Song | null>(null)
  const playlist = ref<Song[]>([])
  const currentIndex = ref(0)
  const isPlaying = ref(false)
  const currentTime = ref(0)
  const volume = ref(80)
  const playMode = ref<'list' | 'random' | 'single'>('list')

  // Getters
  const progress = computed(() => {
    if (!currentSong.value || currentSong.value.duration === 0) return 0
    return (currentTime.value / currentSong.value.duration) * 100
  })

  const hasNext = computed(() => {
    if (playMode.value === 'random') return playlist.value.length > 1
    return currentIndex.value < playlist.value.length - 1
  })

  const hasPrev = computed(() => {
    if (playMode.value === 'random') return playlist.value.length > 1
    return currentIndex.value > 0
  })

  // Actions
  const play = (song: Song) => {
    currentSong.value = song
    isPlaying.value = true
    currentTime.value = 0
  }

  const togglePlay = () => {
    isPlaying.value = !isPlaying.value
  }

  const pause = () => {
    isPlaying.value = false
  }

  const next = () => {
    if (playlist.value.length === 0) return

    if (playMode.value === 'random') {
      let newIndex = Math.floor(Math.random() * playlist.value.length)
      while (newIndex === currentIndex.value && playlist.value.length > 1) {
        newIndex = Math.floor(Math.random() * playlist.value.length)
      }
      currentIndex.value = newIndex
    } else {
      currentIndex.value = (currentIndex.value + 1) % playlist.value.length
    }

    currentSong.value = playlist.value[currentIndex.value]
    isPlaying.value = true
  }

  const prev = () => {
    if (playlist.value.length === 0) return

    if (playMode.value === 'random') {
      next() // In random mode, just pick another random song
      return
    }

    currentIndex.value = (currentIndex.value - 1 + playlist.value.length) % playlist.value.length
    currentSong.value = playlist.value[currentIndex.value]
    isPlaying.value = true
  }

  const setPlaylist = (songs: Song[], index = 0) => {
    playlist.value = songs
    currentIndex.value = index
    if (songs.length > 0) {
      currentSong.value = songs[index]
    }
  }

  const addToPlaylist = (song: Song) => {
    playlist.value.push(song)
  }

  const removeFromPlaylist = (index: number) => {
    playlist.value.splice(index, 1)
    if (index < currentIndex.value) {
      currentIndex.value--
    }
  }

  const setVolume = (val: number) => {
    volume.value = Math.max(0, Math.min(100, val))
  }

  const setCurrentTime = (time: number) => {
    currentTime.value = Math.max(0, Math.min(currentSong.value?.duration || 0, time))
  }

  const setPlayMode = (mode: 'list' | 'random' | 'single') => {
    playMode.value = mode
  }

  return {
    // State
    currentSong,
    playlist,
    currentIndex,
    isPlaying,
    currentTime,
    volume,
    playMode,
    // Getters
    progress,
    hasNext,
    hasPrev,
    // Actions
    play,
    togglePlay,
    pause,
    next,
    prev,
    setPlaylist,
    addToPlaylist,
    removeFromPlaylist,
    setVolume,
    setCurrentTime,
    setPlayMode
  }
})
