<template>
  <n-config-provider :theme="theme">
    <n-layout class="layout">
      <!-- Header -->
      <n-layout-header bordered class="header">
        <div class="header-content">
          <div class="logo">
            <n-icon size="32" color="#18a058">
              <musical-notes />
            </n-icon>
            <span class="logo-text">OpenClaw Music</span>
          </div>
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索音乐、歌手、专辑..."
            class="search-input"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <n-icon><search /></n-icon>
            </template>
          </n-input>
          <n-space>
            <n-button quaternary circle>
              <n-icon><settings /></n-icon>
            </n-button>
          </n-space>
        </div>
      </n-layout-header>

      <n-layout has-sider class="main-layout">
        <!-- Sidebar -->
        <n-layout-sider
          bordered
          collapse-mode="width"
          :collapsed-width="64"
          :width="200"
          :collapsed="collapsed"
          show-trigger
          @collapse="collapsed = true"
          @expand="collapsed = false"
        >
          <n-menu
            :collapsed="collapsed"
            :collapsed-width="64"
            :collapsed-icon-size="22"
            :options="menuOptions"
            :value="activeMenu"
            @update:value="handleMenuSelect"
          />
        </n-layout-sider>

        <!-- Main Content -->
        <n-layout-content class="content">
          <router-view />
        </n-layout-content>
      </n-layout>

      <!-- Player Bar -->
      <n-layout-footer bordered class="player-bar">
        <PlayerBar />
      </n-layout-footer>
    </n-layout>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NConfigProvider,
  NLayout,
  NLayoutHeader,
  NLayoutSider,
  NLayoutContent,
  NLayoutFooter,
  NMenu,
  NInput,
  NButton,
  NIcon,
  NSpace,
  darkTheme,
  type MenuOption
} from 'naive-ui'
import {
  MusicalNotes,
  Search,
  Settings,
  Home,
  List,
  Heart,
  Time,
  Download
} from '@vicons/ionicons5'
import PlayerBar from './components/PlayerBar.vue'

const router = useRouter()
const route = useRoute()

const theme = ref(darkTheme)
const collapsed = ref(false)
const searchKeyword = ref('')
const activeMenu = ref('home')

const menuOptions: MenuOption[] = [
  {
    label: '首页',
    key: 'home',
    icon: () => h(NIcon, null, { default: () => h(Home) })
  },
  {
    label: '歌单',
    key: 'playlist',
    icon: () => h(NIcon, null, { default: () => h(List) })
  },
  {
    label: '收藏',
    key: 'favorites',
    icon: () => h(NIcon, null, { default: () => h(Heart) })
  },
  {
    label: '历史',
    key: 'history',
    icon: () => h(NIcon, null, { default: () => h(Time) })
  },
  {
    label: '下载',
    key: 'downloads',
    icon: () => h(NIcon, null, { default: () => h(Download) })
  }
]

const handleMenuSelect = (key: string) => {
  activeMenu.value = key
  router.push(`/${key}`)
}

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({
      path: '/search',
      query: { keyword: searchKeyword.value }
    })
  }
}
</script>

<style scoped>
.layout {
  height: 100vh;
}

.header {
  height: 64px;
  padding: 0 24px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
  color: #18a058;
}

.search-input {
  width: 400px;
}

.main-layout {
  flex: 1;
  overflow: hidden;
}

.content {
  padding: 24px;
  overflow-y: auto;
}

.player-bar {
  height: 80px;
  padding: 0;
}
</style>
