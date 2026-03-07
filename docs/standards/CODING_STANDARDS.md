# OpenClaw Music Player - Coding Standards

## 1. 通用规范

### 1.1 代码风格

- 使用统一的代码格式化工具
- 保持代码简洁、可读
- 避免过度工程化

### 1.2 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 变量 | 小驼峰 | `userName`, `songList` |
| 常量 | 大写下划线 | `MAX_RETRY_COUNT` |
| 函数 | 小驼峰 | `getUserInfo()` |
| 类/结构体 | 大驼峰 | `UserService`, `SongInfo` |
| 接口 | 大驼峰 + I前缀 | `IMusicParser` |
| 文件 | 小写下划线 | `user_service.go` |

## 2. Go 后端规范

### 2.1 项目结构

```
backend/
├── cmd/                    # 入口程序
│   └── server/
│       └── main.go
├── internal/               # 内部包
│   ├── server/            # HTTP服务
│   ├── parser/            # 音乐解析
│   ├── sync/              # 数据同步
│   └── utils/             # 工具函数
├── pkg/                    # 公共包
├── api/                    # API定义
└── configs/               # 配置文件
```

### 2.2 代码规范

```go
// 包注释
package parser

// MusicParser 音乐解析器接口
type MusicParser interface {
    // Search 搜索音乐
    // keyword: 搜索关键词
    // limit: 返回数量限制
    // offset: 偏移量
    // 返回搜索结果和错误
    Search(keyword string, limit, offset int) (*SearchResult, error)
}

// NetEaseParser 网易云解析器
type NetEaseParser struct {
    client  *http.Client
    baseURL string
}

// NewNetEaseParser 创建网易云解析器
func NewNetEaseParser() *NetEaseParser {
    return &NetEaseParser{
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
        baseURL: "https://music.163.com",
    }
}

// 错误处理
func (p *NetEaseParser) Search(keyword string) (*SearchResult, error) {
    if keyword == "" {
        return nil, fmt.Errorf("keyword cannot be empty")
    }
    
    result, err := p.doSearch(keyword)
    if err != nil {
        return nil, fmt.Errorf("search failed: %w", err)
    }
    
    return result, nil
}
```

### 2.3 错误处理

- 使用 `fmt.Errorf` 包装错误
- 错误信息要小写开头
- 提供上下文信息

```go
// Good
return nil, fmt.Errorf("failed to parse response: %w", err)

// Bad
return nil, errors.New("Failed to parse response")
```

### 2.4 日志规范

```go
logger.Info("Starting server on port %d", port)
logger.Error("Failed to connect to database: %s", err)
logger.Warn("Retry attempt %d/%d", attempt, maxRetries)
```

## 3. TypeScript/Vue 前端规范

### 3.1 项目结构

```
frontend/web/src/
├── components/            # 组件
├── views/                 # 页面
├── stores/                # 状态管理
├── api/                   # API调用
├── composables/           # 组合式函数
├── utils/                 # 工具函数
└── types/                 # 类型定义
```

### 3.2 组件规范

```vue
<template>
  <div class="song-list">
    <!-- 组件内容 -->
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Song } from '../api/music'

// Props 定义
const props = defineProps<{
  songs: Song[]
  loading?: boolean
}>()

// Emits 定义
const emit = defineEmits<{
  play: [song: Song]
  download: [song: Song]
}>()

// 响应式状态
const currentIndex = ref(0)

// 计算属性
const currentSong = computed(() => props.songs[currentIndex.value])

// 方法
const handlePlay = (song: Song) => {
  emit('play', song)
}
</script>

<style scoped>
.song-list {
  /* 样式 */
}
</style>
```

### 3.3 组合式函数规范

```typescript
// useMusic.ts
import { ref } from 'vue'
import * as musicApi from '../api/music'

export function useMusic() {
  const loading = ref(false)
  const error = ref<string | null>(null)

  const search = async (keyword: string) => {
    loading.value = true
    error.value = null
    
    try {
      const result = await musicApi.searchMusic({ keyword })
      return result
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Search failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    search
  }
}
```

### 3.4 类型定义

```typescript
// types/music.ts
export interface Song {
  id: string
  title: string
  artist: string
  album: string
  cover?: string
  duration: number
  url?: string
  source: string
}

export interface Playlist {
  id: string
  name: string
  cover?: string
  songs: Song[]
  createdAt: Date
  updatedAt: Date
}
```

## 4. Git 提交规范

### 4.1 提交信息格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 4.2 Type 类型

| 类型 | 说明 |
|------|------|
| feat | 新功能 |
| fix | 修复bug |
| docs | 文档更新 |
| style | 代码格式调整 |
| refactor | 重构 |
| test | 测试相关 |
| chore | 构建/工具相关 |

### 4.3 示例

```
feat(backend): add NetEase music parser

- Implement search API
- Add URL decryption
- Add lyrics fetch

Closes #123
```

```
fix(frontend): resolve player progress bar issue

The progress bar was not updating correctly when seeking.
Now it syncs with the audio element currentTime.
```

## 5. 测试规范

### 5.1 Go 测试

```go
func TestSearch(t *testing.T) {
    parser := NewNetEaseParser()
    
    result, err := parser.Search("test", 10, 0)
    
    if err != nil {
        t.Errorf("Search failed: %v", err)
    }
    
    if result == nil {
        t.Error("Result should not be nil")
    }
}
```

### 5.2 前端测试

```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SongList from './SongList.vue'

describe('SongList', () => {
  it('renders songs correctly', () => {
    const songs = [
      { id: '1', title: 'Song 1', artist: 'Artist 1' }
    ]
    
    const wrapper = mount(SongList, {
      props: { songs }
    })
    
    expect(wrapper.find('.song-item').exists()).toBe(true)
  })
})
```

## 6. 文档规范

### 6.1 代码注释

- 公共API必须添加注释
- 复杂逻辑需要解释
- 使用中文或英文统一

### 6.2 README 规范

- 项目简介
- 安装说明
- 使用指南
- 贡献指南
- 许可证

## 7. 性能规范

### 7.1 前端性能

- 图片懒加载
- 虚拟滚动（长列表）
- 防抖/节流
- 代码分割

### 7.2 后端性能

- 连接池
- 缓存策略
- 并发控制
- 超时设置
