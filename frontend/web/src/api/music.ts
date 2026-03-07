import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

export interface SearchParams {
  keyword: string
  type?: 'song' | 'album' | 'artist'
  limit?: number
  offset?: number
}

export interface Song {
  id: string
  title: string
  artist: string
  album: string
  cover?: string
  duration: number
  url?: string
  source?: string
}

export interface SearchResult {
  keyword: string
  total: number
  results: Song[]
}

// Search music
export const searchMusic = async (params: SearchParams): Promise<SearchResult> => {
  const response = await api.get('/music/search', { params })
  return response.data
}

// Get song detail
export const getSongDetail = async (id: string): Promise<Song> => {
  const response = await api.get(`/music/detail/${id}`)
  return response.data
}

// Get music URL
export const getMusicUrl = async (id: string, quality = 'standard'): Promise<{ url: string; quality: string }> => {
  const response = await api.get(`/music/url/${id}`, { params: { quality } })
  return response.data
}

// Parse playlist
export const parsePlaylist = async (url: string): Promise<any> => {
  const response = await api.get('/playlist/parse', { params: { url } })
  return response.data
}

// Get lyrics
export const getLyrics = async (id: string): Promise<{ lyrics: string; translated?: string }> => {
  const response = await api.get(`/lyrics/${id}`)
  return response.data
}

// Health check
export const healthCheck = async (): Promise<{ status: string; time: number }> => {
  const response = await api.get('/health')
  return response.data
}

export default {
  searchMusic,
  getSongDetail,
  getMusicUrl,
  parsePlaylist,
  getLyrics,
  healthCheck
}
