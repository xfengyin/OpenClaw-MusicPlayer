# OpenClaw Music Player - API Guide

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **Content-Type**: `application/json`

## 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## API 列表

### 1. 健康检查

**GET** `/health`

**响应**:
```json
{
  "status": "ok",
  "version": "1.0.0",
  "timestamp": 1709836800
}
```

### 2. 搜索音乐

**GET** `/music/search`

**参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| limit | int | 否 | 返回数量，默认20 |
| offset | int | 否 | 偏移量，默认0 |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "keyword": "周杰伦",
    "total": 1000,
    "results": [
      {
        "id": "123456",
        "title": "晴天",
        "artist": "周杰伦",
        "album": "叶惠美",
        "cover": "https://...",
        "duration": 269,
        "source": "netease"
      }
    ]
  }
}
```

### 3. 获取歌曲详情

**GET** `/music/detail/:id`

**参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 歌曲ID |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "123456",
    "title": "晴天",
    "artist": "周杰伦",
    "album": "叶惠美",
    "cover": "https://...",
    "duration": 269
  }
}
```

### 4. 获取播放链接

**GET** `/music/url/:id`

**参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 歌曲ID |
| quality | string | 否 | 音质: standard, medium, high, lossless |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "123456",
    "url": "https://...",
    "quality": "high"
  }
}
```

### 5. 获取歌词

**GET** `/lyrics/:id`

**参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 歌曲ID |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "123456",
    "lyrics": "[00:00.00]歌词内容...",
    "translated": "[00:00.00]翻译内容..."
  }
}
```

### 6. 解析歌单

**GET** `/playlist/parse`

**参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| url | string | 是 | 歌单链接 |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "url": "https://music.163.com/playlist/...",
    "title": "歌单名称",
    "songs": [
      {
        "id": "123456",
        "title": "歌曲名",
        "artist": "艺术家"
      }
    ]
  }
}
```

### 7. 创建下载任务

**POST** `/download`

**请求体**:
```json
{
  "id": "123456",
  "quality": "high"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "taskId": "task_123456",
    "status": "pending"
  }
}
```

### 8. 获取下载进度

**GET** `/download/progress/:taskId`

**参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| taskId | string | 是 | 任务ID |

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "taskId": "task_123456",
    "status": "downloading",
    "progress": 45.5,
    "speed": "1.5MB/s",
    "totalSize": 10240000
  }
}
```

## 使用示例

### 搜索歌曲

```bash
curl "http://localhost:8080/api/v1/music/search?keyword=周杰伦&limit=10"
```

### 获取播放链接

```bash
curl "http://localhost:8080/api/v1/music/url/123456?quality=high"
```

### 下载歌曲

```bash
curl -X POST "http://localhost:8080/api/v1/download" \
  -H "Content-Type: application/json" \
  -d '{"id": "123456", "quality": "high"}'
```
