<template>
  <div class="search-bar">
    <n-input-group>
      <n-input
        v-model:value="keyword"
        :placeholder="placeholder"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <n-icon><search /></n-icon>
        </template>
      </n-input>
      <n-button type="primary" @click="handleSearch">
        搜索
      </n-button>
    </n-input-group>

    <!-- Search History -->
    <div v-if="showHistory && searchHistory.length > 0" class="search-history">
      <div class="history-header">
        <span>搜索历史</span>
        <n-button text size="tiny" @click="clearHistory">
          清空
        </n-button>
      </div>
      <n-space>
        <n-tag
          v-for="item in searchHistory"
          :key="item"
          closable
          @click="keyword = item; handleSearch()"
          @close="removeHistory(item)"
        >
          {{ item }}
        </n-tag>
      </n-space>
    </div>

    <!-- Hot Search -->
    <div v-if="showHot" class="hot-search">
      <div class="hot-header">热门搜索</div>
      <n-space>
        <n-tag
          v-for="(item, index) in hotSearches"
          :key="index"
          :type="index < 3 ? 'error' : 'default'"
          @click="keyword = item; handleSearch()"
        >
          {{ index + 1 }}. {{ item }}
        </n-tag>
      </n-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NInput, NInputGroup, NButton, NIcon, NTag, NSpace } from 'naive-ui'
import { Search } from '@vicons/ionicons5'

const props = defineProps<{
  placeholder?: string
  showHistory?: boolean
  showHot?: boolean
}>()

const emit = defineEmits<{
  search: [keyword: string]
}>()

const keyword = ref('')
const searchHistory = ref<string[]>([])
const hotSearches = ref([
  '周杰伦',
  '林俊杰',
  '薛之谦',
  '邓紫棋',
  '陈奕迅',
  'Taylor Swift',
  'Ed Sheeran',
  '流行歌曲',
  '经典老歌',
  '摇滚'
])

onMounted(() => {
  // Load search history from localStorage
  const history = localStorage.getItem('searchHistory')
  if (history) {
    searchHistory.value = JSON.parse(history)
  }
})

const handleSearch = () => {
  if (!keyword.value.trim()) return
  
  // Add to history
  addToHistory(keyword.value)
  
  emit('search', keyword.value)
}

const addToHistory = (term: string) => {
  // Remove if exists
  const index = searchHistory.value.indexOf(term)
  if (index > -1) {
    searchHistory.value.splice(index, 1)
  }
  
  // Add to front
  searchHistory.value.unshift(term)
  
  // Keep only 10 items
  if (searchHistory.value.length > 10) {
    searchHistory.value.pop()
  }
  
  // Save to localStorage
  localStorage.setItem('searchHistory', JSON.stringify(searchHistory.value))
}

const removeHistory = (term: string) => {
  const index = searchHistory.value.indexOf(term)
  if (index > -1) {
    searchHistory.value.splice(index, 1)
    localStorage.setItem('searchHistory', JSON.stringify(searchHistory.value))
  }
}

const clearHistory = () => {
  searchHistory.value = []
  localStorage.removeItem('searchHistory')
}
</script>

<style scoped>
.search-bar {
  width: 100%;
}

.search-history,
.hot-search {
  margin-top: 16px;
}

.history-header,
.hot-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.hot-header {
  justify-content: flex-start;
}
</style>
