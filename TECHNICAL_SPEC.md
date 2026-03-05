# OpenClaw Music Player - 技术规范

## API 规范

### 响应格式
```json
{
  "code": 200,
  "msg": "success",
  "data": {}
}
```

### 错误码
- 10001: 认证失败
- 20001: 下载失败
- 30001: 解析失败

## Go 代码规范

- 使用 gofmt 格式化
- 函数首字母大写为导出
- 错误显式返回
- 并发安全使用 mutex/channel

## 前端代码规范

- 目录: kebab-case
- 组件: PascalCase
- Hooks: camelCase
- 使用 TypeScript 类型
