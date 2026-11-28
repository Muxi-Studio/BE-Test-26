# blog-backend

木犀博客后端的待办接口，数据存在内存里，重启即丢。

## 开发

```bash
go run .         # 监听 :8080
go test ./...    # 跑测试
```

## 目录

- `internal/api`：HTTP 接口与路由
- `internal/store`：内存存储
- `docs/`：接口说明、事故记录
