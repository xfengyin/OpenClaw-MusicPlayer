# OpenClaw Music Player - System Design

## 1. 系统架构

### 1.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      Client Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Web App    │  │ Electron App │  │   Mobile     │      │
│  │  (Vue3)      │  │  (Desktop)   │  │   (Future)   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTP/WebSocket
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway Layer                        │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Go HTTP Server (Gin)                   │   │
│  │  • RESTful API  • CORS  • Rate Limit  • Auth       │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Service Layer                             │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐           │
│  │   Parser   │  │  Download  │  │   Lyrics   │           │
│  │  Service   │  │  Service   │  │  Service   │           │
│  └────────────┘  └────────────┘  └────────────┘           │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  Data Source Layer                          │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐           │
│  │  NetEase   │  │    QQ      │  │   Local    │           │
│  │   Music    │  │   Music    │  │   Storage  │           │
│  └────────────┘  └────────────┘  └────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 前端 | Vue3 + TypeScript | 响应式UI框架 |
| 桌面 | Electron | 跨平台桌面应用 |
| 后端 | Go + Gin | 高性能HTTP服务 |
| 数据库 | SQLite | 本地数据存储 |
| 缓存 | In-Memory LRU | 内存缓存 |
| 构建 | Vite | 前端构建工具 |

## 2. 核心模块设计

### 2.1 音乐解析模块 (Parser)

```go
type MusicParser interface {
    Search(keyword string, limit, offset int) (*SearchResult, error)
    GetSongURL(id string, quality int) (string, error)
    GetLyrics(id string) (string, string, error)
}
```

**职责**:
- 对接第三方音乐平台API
- 歌曲搜索
- 获取播放链接
- 歌词获取

**实现**:
- NetEaseParser: 网易云音乐
- QQParser: QQ音乐 (未来)
- LocalParser: 本地音乐

### 2.2 播放控制模块 (Player)

```typescript
interface PlayerState {
  currentSong: Song | null
  isPlaying: boolean
  currentTime: number
  volume: number
  playMode: 'list' | 'random' | 'single'
  playlist: Song[]
}
```

**职责**:
- 播放/暂停/停止
- 进度控制
- 音量控制
- 播放模式切换
- 播放列表管理

### 2.3 下载管理模块 (Download)

```go
type DownloadTask struct {
    ID        string
    SongID    string
    URL       string
    Status    string // pending, downloading, completed, failed
    Progress  float64
    Speed     string
    LocalPath string
}
```

**职责**:
- 多线程下载
- 断点续传
- 下载队列管理
- 进度通知

## 3. API 设计

### 3.1 RESTful API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/health | 健康检查 |
| GET | /api/v1/music/search | 搜索音乐 |
| GET | /api/v1/music/detail/:id | 歌曲详情 |
| GET | /api/v1/music/url/:id | 获取播放链接 |
| GET | /api/v1/lyrics/:id | 获取歌词 |
| POST | /api/v1/download | 创建下载任务 |
| GET | /api/v1/download/progress/:id | 下载进度 |

### 3.2 WebSocket API

用于实时推送:
- 下载进度
- 播放状态
- 歌词同步

## 4. 数据模型

### 4.1 Song (歌曲)

```typescript
interface Song {
  id: string           // 歌曲ID
  title: string        // 标题
  artist: string       // 艺术家
  album: string        // 专辑
  cover?: string       // 封面URL
  duration: number     // 时长(秒)
  url?: string         // 播放URL
  source: string       // 来源
}
```

### 4.2 Playlist (歌单)

```typescript
interface Playlist {
  id: string
  name: string
  cover?: string
  songs: Song[]
  createdAt: Date
  updatedAt: Date
}
```

### 4.3 DownloadTask (下载任务)

```typescript
interface DownloadTask {
  id: string
  song: Song
  status: 'pending' | 'downloading' | 'completed' | 'failed'
  progress: number
  speed: string
  totalSize: number
  downloadedSize: number
}
```

## 5. 缓存策略

### 5.1 内存缓存 (LRU)

```go
type Cache struct {
    data map[string]cacheItem
    ttl  time.Duration
}
```

**缓存内容**:
- 搜索结果 (30分钟)
- 歌曲URL (10分钟)
- 歌词 (1小时)

### 5.2 本地存储

**SQLite 表结构**:
- `songs`: 歌曲信息
- `playlists`: 歌单
- `play_history`: 播放历史
- `download_tasks`: 下载任务

## 6. 安全设计

### 6.1 CORS

```go
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        // ...
    }
}
```

### 6.2 限流

```go
rateLimit := ratelimit.NewLimiter(100, time.Minute)
```

## 7. 性能优化

### 7.1 前端优化

- 虚拟滚动 (长列表)
- 图片懒加载
- 组件按需加载
- 状态管理优化

### 7.2 后端优化

- 连接池
- 并发控制
- 缓存策略
- 压缩传输

## 8. 错误处理

######