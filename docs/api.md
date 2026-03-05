# OpenClaw Music Player API 文档

## 基础信息

- **Base URL**: /api/v1
- **Accept**: application/json
- **Content-Type**: application/json

## 认证

部分接口需要认证，.Token 放在请求头中：
```
Authorization: Bearer {token}
```

## 接口列表

### 健康检查

`GET /api/v1/health`

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "status": "ok"
  }
}
```

### 搜索音乐

`GET /api/v1/parser/search?keyword=xxx&page=1&size=20`

### 获取音乐详情

`GET /api/v1/parser/detail/:id`

### 获取音乐 URL

`GET /api/v1/parser/url/:id`

### 开始下载

`POST /api/v1/download/start`

```json
{
  "music_id": "123456",
  "title": "歌曲名",
  "url": "https://..."
}
```

### 获取下载状态

`GET /api/v1/download/status/:id`

### 播放控制

`POST /api/v1/playback/play`

```json
{
  "music_id": "123456"
}
```

### 远程控制

`POST /api/v1/remote/control`

```json
{
  "command": "play",
  "data": {}
}
```
