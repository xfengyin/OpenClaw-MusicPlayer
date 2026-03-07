<template>
  <div class="search-page">
    <SearchBar
      placeholder="搜索音乐、歌手、专辑..."
      :show-history="true"
      :show-hot="true"
      @search="onSearch"
    />

    <!-- Search Results -->
    <div v-if="loading" class="loading">
      <n-spin size="large" />
    </div>

    <div v-else-if="error" class="error">
      <n-alert type="error" :title="error" />
    </div>

    <div v-else-if="hasSearched" class="results">
      <n-tabs v-model:value="activeTab" type="line">
        <n-tab-pane name="songs" :tab="`歌曲 ${results.length}`">
          <SongList
            :songs="results"
            @play="onPlay"
            @download="onDownload"
          />
        </n-tab-pane>
        
        <n-tab-pane name="playlists" tab="歌单">
          <n-empty description="歌单搜索功能开发中..." />
        </n-tab-pane>
        
        <n-tab-pane name="artists" tab="歌手">
          <n-empty description="歌手搜索功能开发中..." />
        </n-tab-pane>
      </n-tabs>
    </div>

    <div v-else class="empty">
      <n-empty description="输入关键词开始搜索" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NSpin, NAlert, NEmpty, NTabs, NTabPane } from 'naive-ui'
import SearchBar from '../components/SearchBar.vue'
import SongList from '../components/SongList.vue'
import { useMusic } from '../composables/useMusic'
import { usePlayerStore } from '../stores/player'
import type { Song } from '../api/music'

const route = useRoute()
const router = useRouter()
const playerStore = usePlayerStore()
const { search, loading, error } = useMusic()

const keyword = ref('')
const activeTab = ref('songs')
const results = ref<Song[]>([])
const hasSearched = ref(false)

// Get keyword from URL
const urlKeyword = computed(() => route.query.keyword as string)

// Search on mount if keyword exists
if (urlKeyword.value) {
  keyword.value = urlKeyword.value
  performSearch(urlKeyword.value)
}

const onSearch = async (term: string) => {
  keyword.value = term
  
  // Update URL
  router.push({
    path: '/search',
    query: { keyword: term }
  })
  
  await performSearch(term)
}

const performSearch = async (term: string) => {
  hasSearched.value = true
  
  try {
    const response = await search({
      keyword: term,
      limit: 30,
      offset: 0
    })
    
    results.value = response.results
  } catch (err) {
    console.error('Search error:', err)
    results.value = []
  }
}

const onPlay = async (song: Song) => {
  try {
    // Get song URL if not present
    if (!song.url) {
      const { getSongUrl } = useMusic()
      song.url = await getSongUrl(song.id)
    }
    
    // Play song
    playerStore.play(song)
  } catch (err) {
    console.error('Failed to play song:', err)
  }
}

const onDownload = (song: Song) => {
  // TODO: Implement download
  console.log('Download:', song)
}
</script>

<style scoped>
.search-page {
  padding: 24px;
}

.loading,
.error,
.empty {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px;
}

.results {
  margin-top: 24px;
}
</style>
