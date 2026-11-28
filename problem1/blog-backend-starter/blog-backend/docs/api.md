# 接口说明

所有响应都套一层信封：

```json
{ "code": 0, "msg": "ok", "data": null }
```

## 列表

`GET /todos`

`data` 是待办数组。

```json
{
  "code": 0,
  "msg": "ok",
  "data": [{ "id": 1, "title": "写题解", "finished": false }]
}
```

## 错误

失败时 `code` 非 0，`msg` 是给人看的说明。

## 字段

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | int | 自增 ID |
| `title` | string | 标题 |
| `finished` | bool | 是否完成 |

## 错误码

| `code` | 含义 |
| --- | --- |
| 0 | 成功 |
| 1 | 服务端错误 |
