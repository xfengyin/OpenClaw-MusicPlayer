# OpenClaw Music Player

[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.19%2B-blue)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4-blue)](https://vuejs.org/)
[![Electron](https://img.shields.io/badge/Electron-25%2B-47848F)](https://www.electronjs.org/)

全功能音乐播放器，对标 AlgerMusicPlayer，基于 Go + Vue3 + Electron 构建。

## 🚀 特性

- **网易云音乐生态** - 登录、歌单、下载、歌词同步
- **高性能服务层** - Go 语言核心服务，支持并发下载
- **跨平台支持** - Windows/macOS/Linux 桌面端
- **Web 端支持** - 独立部署，远程控制
- **现代 UI** - naive-ui + Tailwind CSS
- **开发者友好** - 开源，社区驱动

## 📊 项目结构

```
OpenClaw-MusicPlayer/
├── cmd/                  # 可执行文件入口
│   ├── server/          # Go 服务
│   └── player/          # Electron 启动
├── server/              # Go 服务层
│   ├── handler/         # HTTP 处理器
│   ├── middleware/      # 中间件
│   ├── service/         # 核心服务
│   └── utils/           # 工具函数
├── web/                 # 前端项目
│   ├── src/
│   └── public/
├── electron/            # Electron 配置
│   ├── main/
│   ├── renderer/
│   └── build/
└── docs/                # 文档
```

## 🛠️ 技术栈

### 后端/服务层
- **框架**: Gin, gRPC
- **数据库**: BoltDB, SQLite
- **并发**: goroutine + channel

### 前端
- **框架**: Vue3 + TypeScript
- **UI**: naive-ui + Tailwind CSS
- **状态**: Pinia

### 桌面端
- **封装**: Electron 25+

## 🚀 快速开始

```bash
# 1. 克隆项目
git clone https://github.com/xfengyin/OpenClaw-MusicPlayer.git
cd OpenClaw-MusicPlayer

# 2. 启动 Go 服务
cd server
go run cmd/server/main.go

# 3. 启动前端
cd ../web
npm install
npm run dev
```

## 📖 API 文档

Go 服务提供 RESTful API：

- `GET /api/v1/health` - 健康检查
- `GET /api/v1/parser/search?keyword=xxx` - 搜索音乐
- `POST /api/v1/download/start` - 开始下载

## 📄 许可证

[MIT](LICENSE)
