<template>
  <div class="home">
    <n-h1>发现音乐</n-h1>
    
    <!-- Featured Playlists -->
    <n-h2>推荐歌单</n-h2>
    <n-grid :cols="5" :x-gap="16" :y-gap="16">
      <n-grid-item v-for="i in 10" :key="i">
        <PlaylistCard
          :title="`推荐歌单 ${i}`"
          :cover="`https://picsum.photos/200/200?random=${i}`"
          :play-count="10000 + i * 1000"
        />
      </n-grid-item>
    </n-grid>

    <!-- New Songs -->
    <n-h2>新歌速递</n-h2>
    <SongList :songs="newSongs" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NH1, NH2, NGrid, NGridItem } from 'naive-ui'
import PlaylistCard from '../components/PlaylistCard.vue'
import SongList from '../components/SongList.vue'

const newSongs = ref([
  { id: '1', title: '歌曲 1', artist: '歌手 A', album: '专辑 X', duration: 240 },
  { id: '2', title: '歌曲 2', artist: '歌手 B', album: '专辑 Y', duration: 180 },
  { id: '3', title: '歌曲 3', artist: '歌手 C', album: '专辑 Z', duration: 210 },
  { id: '4', title: '歌曲 4', artist: '歌手 D', album: '专辑 W', duration: 195 },
  { id: '5', title: '歌曲 5', artist: '歌手 E', album: '专辑 V', duration: 225 },
])
</script>

<style scoped>
.home {
  padding-bottom: 24px;
}

.n-h1 {
  margin-bottom: 24px;
}

.n-h2 {
  margin-top: 32px;
  margin-bottom: 16px;
}
</style>
