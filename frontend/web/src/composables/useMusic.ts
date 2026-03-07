import { ref } from 'vue'
import * as musicApi from '../api/music'
import type { Song, SearchParams } from '../api/music'

export function useMusic() {
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Search music
  const search = async (params: SearchParams) => {
    loading.value = true
    error.value = null
    
    try {
      const result = await musicApi.searchMusic(params)
      return result
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Search failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Get song URL
  const getSongUrl = async (id: string, quality?: string) => {
    loading.value = true
    error.value = null
    
    try {
      const result = await musicApi.getMusicUrl(id, quality)
      return result.url
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to get URL'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Get lyrics
  const getLyrics = async (id: string) => {
    loading.value = true
    error.value = null
    
    try {
      const result = await musicApi.getLyrics(id)
      return result
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to get lyrics'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Parse playlist
  const parsePlaylist = async (url: string) => {
    loading.value = true
    error.value = null
    
    try {
      const result = await musicApi.parsePlaylist(url)
      return result
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to parse playlist'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    search,
    getSongUrl,
    getLyrics,
    parsePlaylist
  }
}
