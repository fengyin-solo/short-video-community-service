# 短视频平台（short-video）

纯 Go 标准库实现的短视频平台后端，覆盖用户、分类、视频、标签、点赞、评论、关注、收藏、举报、通知、播单等完整社交互动闭环。零第三方依赖，开箱即跑。

## 技术栈

- Go 1.22+，仅使用标准库（`net/http` + 标准库）
- 标准工程分层：`cmd` / `internal`（app/config/model/store/service/handler）/ `pkg`
- 内存存储（`store.MemoryStore`），线程安全

## 运行

```bash
cd origin
go run ./cmd/server
# 或
go build -o server ./cmd/server && ./server
```

环境变量：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | 8080 | 监听端口 |
| `ADDR` | `:8080` | 监听地址（优先于 PORT） |
| `MAX_PAGE_SIZE` | 100 | 分页最大页大小 |
| `AUTH_TOKEN` | `dev-token` | API 鉴权令牌（`X-Auth-Token` 请求头） |
| `RATE_LIMIT` | 600 | 每 IP 每分钟限流次数 |
| `LOG_LEVEL` | info | 日志级别 debug/info/warn/error |

服务启动后可先调用 `POST /api/seed/demo` 初始化演示数据。

## 统一响应

```json
{"code":0,"message":"ok","data":{...}}
```

错误映射：`ValidationError→400`、`ErrNotFound→404`、`ErrConflict→409`、其他→500。

## API 一览

鉴权：所有 `/api/` 接口需在请求头携带 `X-Auth-Token`。

### 用户 /api/users
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/users | 创建用户 |
| GET | /api/users | 列表（status/keyword 筛选） |
| GET | /api/users/{id} | 详情 |
| PUT | /api/users/{id} | 更新 |
| DELETE | /api/users/{id} | 删除 |
| POST | /api/users/{id}/ban | 封禁 |
| POST | /api/users/{id}/unban | 解封 |
| GET | /api/users/{id}/videos | 该用户的视频 |

### 分类 /api/categories
| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/GET{id}/PUT/DELETE | /api/categories | 分类 CRUD |

### 视频 /api/videos
状态机：`draft→reviewing→published/banned`，published 可 `banned`。
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/videos | 创建视频（草稿） |
| GET | /api/videos | 列表（status/category_id/author_id/keyword） |
| GET | /api/videos/{id} | 详情 |
| DELETE | /api/videos/{id} | 删除 |
| POST | /api/videos/{id}/submit | 提交审核 |
| POST | /api/videos/{id}/approve | 审核通过（发布） |
| POST | /api/videos/{id}/ban | 下架 |
| POST | /api/videos/{id}/view | 播放量 +1 |

### 标签 /api/tags
| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/GET{id}/DELETE | /api/tags | 标签管理 |

### 点赞 /api/likes
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/likes | 点赞（body: user_id, video_id） |
| DELETE | /api/likes | 取消点赞 |
| GET | /api/likes | 列表 |

### 评论 /api/comments
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/comments | 发评论（支持 parent_id 楼中楼） |
| GET | /api/comments | 列表 |
| GET | /api/comments/{id} | 详情 |
| DELETE | /api/comments/{id} | 删除（软删） |

### 关注 /api/follows
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/follows | 关注 |
| DELETE | /api/follows | 取消关注 |
| GET | /api/follows | 列表 |

### 收藏 /api/favorites
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/favorites | 收藏 |
| DELETE | /api/favorites | 取消收藏 |
| GET | /api/favorites | 列表 |

### 举报 /api/reports
状态机：`pending→processed/rejected`；处理（process）可联动下架视频。
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/reports | 创建举报 |
| GET | /api/reports | 列表 |
| GET | /api/reports/{id} | 详情 |
| POST | /api/reports/{id}/process | 处理（有效） |
| POST | /api/reports/{id}/reject | 驳回 |
| DELETE | /api/reports/{id} | 删除 |

### 通知 /api/notifications
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/notifications | 创建通知 |
| GET | /api/notifications | 列表（user_id/type/read） |
| GET | /api/notifications/{id} | 详情 |
| POST | /api/notifications/{id}/read | 标记已读 |
| POST | /api/notifications/read-all | 全部已读 |
| DELETE | /api/notifications/{id} | 删除 |

### 播单 /api/playlists
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/playlists | 创建播单 |
| GET | /api/playlists | 列表 |
| GET | /api/playlists/{id} | 详情 |
| PUT | /api/playlists/{id} | 更新 |
| DELETE | /api/playlists/{id} | 删除 |
| POST | /api/playlists/{id}/videos | 添加视频 |
| DELETE | /api/playlists/{id}/videos/{videoID} | 移除视频 |

### 视频标签 /api/video-tags
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/video-tags | 打标签 |
| DELETE | /api/video-tags | 取消标签 |
| GET | /api/video-tags | 列表 |
| GET | /api/videos/{id}/tags | 视频的标签名列表 |

### 搜索 / 信息流
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/search?q=关键词 | 搜索已发布视频（相关度排序） |
| GET | /api/feed | 信息流（按播放量） |

### 其他
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/seed/demo | 初始化演示数据 |
