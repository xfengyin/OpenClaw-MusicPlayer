# OpenClaw Music Player 🎵

一款开源的跨平台音乐播放器，支持多音源搜索、在线播放和下载管理。

[![CI](https://github.com/xfengyin/OpenClaw-MusicPlayer/actions/workflows/team-ci.yml/badge.svg)](https://github.com/xfengyin/OpenClaw-MusicPlayer/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 特性

- 🎵 **多音源支持**: 网易云音乐、QQ音乐等
- 🔍 **智能搜索**: 歌曲、歌手、专辑搜索
- 📥 **下载管理**: 多线程下载，支持断点续传
- 🎨 **精美界面**: 现代化UI设计，支持暗黑模式
- 💻 **跨平台**: Web、桌面(Windows/macOS/Linux)
- 📝 **歌词显示**: 实时歌词同步
- 📱 **响应式设计**: 适配各种屏幕尺寸

## 技术栈

### 前端
- Vue 3 + TypeScript
- Vite
- Naive UI
- Pinia

### 后端
- Go
- Gin Framework
- SQLite

### 桌面
- Electron

## 快速开始

### 环境要求

- Node.js 20+
- Go 1.21+
- Git

### 安装

```bash
# 克隆仓库
git clone https://github.com/xfengyin/OpenClaw-MusicPlayer.git
cd OpenClaw-MusicPlayer

# 安装前端依赖
cd frontend/web
npm install

# 安装后端依赖
cd ../../backend
go mod download

# 安装 Electron 依赖
cd ../frontend/electron
npm install
```

### 开发

```bash
# 启动后端服务
cd backend
go run cmd/server/main.go

# 启动前端开发服务器 (新终端)
cd frontend/web
npm run dev

# 启动 Electron 桌面应用 (新终端)
cd frontend/electron
npm run dev
```

### 构建

```bash
# 构建后端
./scripts/build/build-backend.sh

# 构建前端
./scripts/build/build-frontend.sh

# 构建桌面应用
./scripts/build/build-electron.sh
```

## 项目结构

```
OpenClaw-MusicPlayer/
├── backend/              # Go 后端
│   ├── cmd/             # 入口程序
│   ├── internal/        # 内部模块
│   ├── api/             # API 文档
│   └── configs/         # 配置文件
├── frontend/            # 前端
│   ├── web/            # Vue3 Web 应用
│   └── electron/       # Electron 桌面应用
├── tests/              # 测试
├── docs/               # 文档
│   ├── architecture/   # 架构文档
│   ├── api/            # API 文档
│   └── standards/      # 编码规范
└── scripts/            # 构建脚本
```

## 团队开发

### 分支策略

- `main`: 生产分支
- `team-dev`: 团队开发主干
- `feature/*`: 功能分支

### 团队分工

| 角色 | 人数 | 负责模块 |
|------|------|----------|
| 技术总监 | 1 | 架构设计、团队协调 |
| Go 后端 | 2 | backend/ 目录 |
| 前端 | 2 | frontend/ 目录 |
| 测试 | 1 | tests/ 目录 |
| DevOps | 1 | CI/CD、部署 |

### 开发流程

1. 从 `team-dev` 创建功能分支
2. 开发完成后提交 PR 到 `team-dev`
3. Code Review 后合并
4. 定期将 `team-dev` 合并到 `main`

## API 文档

详见 [API Guide](docs/api/API_GUIDE.md)

## 架构设计

详见 [System Design](docs/architecture/SYSTEM_DESIGN.md)

## 编码规范

详见 [Coding Standards](docs/standards/CODING_STANDARDS.md)

## 贡献指南

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

[MIT](LICENSE) License

## 致谢

感谢所有贡献者的支持！
