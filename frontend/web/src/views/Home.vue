<template>
  <div class="home">
    <!-- Banner -->
    <div class="banner">
      <n-carousel autoplay>
        <div
          v-for="i in 5"
          :key="i"
          class="banner-item"
          :style="{ backgroundImage: `url(https://picsum.photos/1200/300?random=${i})` }"
        >
          <div class="banner-content">
            <h2>发现好音乐 {{ i }}</h2>
            <p>海量曲库，随心畅听</p>
          </div>
        </div>
      </n-carousel>
    </div>

    <!-- Featured Playlists -->
    <section class="section">
      <div class="section-header">
        <h2>推荐歌单</h2>
        <n-button text @click="$router.push('/playlist')">
          查看更多
          <n-icon><chevron-forward /></n-icon>
        </n-button>
      </div>
      <n-grid :cols="5" :x-gap="16" :y-gap="16">
        <n-grid-item v-for="i in 10" :key="i">
          <PlaylistCard
            :title="`推荐歌单 ${i}`"
            :cover="`https://picsum.photos/300/300?random=${i}`"
            :play-count="10000 + i * 1000"
            @click="$router.push(`/playlist/${i}`)"
          />
        </n-grid-item>
      </n-grid>
    </section>

    <!-- New Songs -->
    <section class="section">
      <div class="section-header">
        <h2>新歌速递</h2>
        <n-button text @click="$router.push('/search')">
          查看更多
          <n-icon><chevron-forward /></n-icon>
        </n-button>
      </div>
      <SongList :songs="newSongs" @play="onPlay" />
    </section>

    <!-- Hot Artists -->
    <section class="section">
      <div class="section-header">
        <h2>热门歌手</h2>
        <n-button text>
          查看更多
          <n-icon><chevron-forward /></n-icon>
        </n-button>
      </div>
      <n-grid :cols="6" :x-gap="16" :y-gap="16">
        <n-grid-item v-for="i in 12" :key="i">
          <div class="artist-card" @click="$router.push(`/search?keyword=歌手${i}`)">
            <n-avatar
              :size="100"
              :src="`https://picsum.photos/200/200?random=${i + 100}`"
              round
            />
            <div class="artist-name">歌手 {{ i }}</div>
          </div>
        </n-grid-item>
      </n-grid>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NCarousel, NGrid, NGridItem, NButton, NIcon, NAvatar } from 'naive-ui'
import { ChevronForward } from '@vicons/ionicons5'
import PlaylistCard from '../components/PlaylistCard.vue'
import SongList from '../components/SongList.vue'
import { usePlayerStore } from '../stores/player'
import type { Song } from '../api/music'

const playerStore = usePlayerStore()

const newSongs = ref<Song[]>([
  { id: '1', title: '晴天', artist: '周杰伦', album: '叶惠美', duration: 269 },
  { id: '2', title: '七里香', artist: '周杰伦', album: '七里香', duration: 299 },
  { id: '3', title: '稻香', artist: '周杰伦', album: '魔杰座', duration: 223 },
  { id: '4', title: '夜曲', artist: '周杰伦', album: '十一月的萧邦', duration: 226 },
  { id: '5', title: '告白气球', artist: '周杰伦', album: '周杰伦的床边故事', duration: 215 },
])

const onPlay = (song: Song) => {
  playerStore.play(song)
}
</script>

<style scoped>
.home {
  padding-bottom: 24px;
}

.banner {
  margin-bottom: 32px;
  border-radius: 12px;
  overflow: hidden;
}

.banner-item {
  height: 300px;
  background-size: cover;
  background-position: center;
  display: flex;
  align-items: flex-end;
  padding: 32px;
}

.banner-content {
  color: white;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
}

.banner-content h2 {
  font-size: 32px;
  margin-bottom: 8px;
}

.banner-content p {
  font-size: 16px;
  opacity: 0.9;
}

.section {
  margin-bottom: 32px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h2 {
  font-size: 20px;
  font-weight: bold;
}

.artist-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s;
}

.artist-card:hover {
  transform: translateY(-4px);
}

.artist-name {
  margin-top: 8px;
  font-size: 14px;
  text-align: center;
}
</style>
