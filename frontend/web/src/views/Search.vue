<template>
  <div class="search">
    <n-h1>搜索结果: {{ keyword }}</n-h1>
    
    <n-tabs v-model:value="activeTab" type="line">
      <n-tab-pane name="songs" tab="歌曲">
        <SongList :songs="searchResults" />
      </n-tab-pane>
      <n-tab-pane name="artists" tab="歌手">
        <n-empty description="暂无歌手数据" />
      </n-tab-pane>
      <n-tab-pane name="albums" tab="专辑">
        <n-empty description="暂无专辑数据" />
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { NH1, NTabs, NTabPane, NEmpty } from 'naive-ui'
import SongList from '../components/SongList.vue'

const route = useRoute()
const keyword = ref('')
const activeTab = ref('songs')

const searchResults = ref([
  { id: '1', title: '搜索结果 1', artist: '歌手 A', album: '专辑 X', duration: 240 },
  { id: '2', title: '搜索结果 2', artist: '歌手 B', album: '专辑 Y', duration: 180 },
  { id: '3', title: '搜索结果 3', artist: '歌手 C', album: '专辑 Z', duration: 210 },
])

watch(() => route.query.keyword, (newKeyword) => {
  if (newKeyword) {
    keyword.value = newKeyword as string
    // TODO: Call API to search
  }
}, { immediate: true })
</script>

<style scoped>
.search {
  padding-bottom: 24px;
}

.n-h1 {
  margin-bottom: 24px;
}
</style>
