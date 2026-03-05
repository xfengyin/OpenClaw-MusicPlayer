# OpenClaw Music Player

[![CI](https://github.com/xfengyin/OpenClaw-MusicPlayer/actions/workflows/ci.yml/badge.svg)](https://github.com/xfengyin/OpenClaw-MusicPlayer/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

全功能音乐播放器，对标 AlgerMusicPlayer，基于 Go + Vue3 + Electron 构建。

## 🚀 功能特性

- ✅ 网易云音乐登录/同步
- ✅ 歌单/专辑/MV 浏览
- ✅ 音乐下载 (多线程)
- ✅ 歌词显示与转换
- ✅ 远程控制
- ✅ 桌面/Web/移动端支持

## 📊 项目结构

```
OpenClaw-MusicPlayer/
├── cmd/              # 可执行文件入口
├── server/           # Go 服务层
├── web/              # Vue3 前端
├── electron/         # Electron 配置
└── docs/             # 文档
```

## 🛠️ 快速开始

```bash
# 安装服务端依赖
cd server
go mod download

# 安装前端依赖
cd ../web
npm install

# 启动服务
cd ../server
go run cmd/server/main.go

# 启动前端 (新窗口)
cd ../web
npm run dev
```

## 📖 文档

- [API 文档](docs/api.md)
- [部署文档](docs/deployment.md)
- [项目计划](PROJECT_PLAN.md)
- [技术规范](TECHNICAL_SPEC.md)

## 📄 许可证

[MIT](LICENSE)

## 🙏 致谢

- 核心功能参考: AlgerMusicPlayer
- 网易云 API: 开源实现参考
