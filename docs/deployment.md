# OpenClaw Music Player - 部署文档

## 环境要求

- Go 1.19+
- Node.js 20+
- React Native CLI (可选)

## 本地开发

### 1. 克隆项目

```bash
git clone https://github.com/xfengyin/OpenClaw-MusicPlayer.git
cd OpenClaw-MusicPlayer
```

### 2. 启动服务端

```bash
cd server
go run cmd/server/main.go
```

### 3. 启动前端

```bash
cd web
npm install
npm run dev
```

## 生产部署

### Web 端部署

1. 构建前端
```bash
cd web
npm run build
```

2. 配置 Nginx
```nginx
location /api {
    proxy_pass http://localhost:8080;
}
location / {
    root /path/to/dist;
    index index.html;
}
```

### 桌面端打包

```bash
cd electron
npm install
npm run build
```

## Docker 部署

```bash
docker build -t openclaw-server .
docker run -p 8080:8080 openclaw-server
```
