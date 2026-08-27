# 排行榜事件处理服务

纯 Go 标准库实现的排行榜事件处理 API 服务，负责赛季、成员、得分事件、排名状态和变动历史的生命周期管理，零第三方依赖。

## 运行说明

在 `origin/` 目录下执行：

```bash
go run ./cmd/server
```

默认监听 `:8080`，所有能力通过 HTTP API 提供。

环境变量：
- `PORT` — 服务端口，默认 8080
- `ADDR` — 监听地址（覆盖 PORT）
- `MAX_PAGE_SIZE` — 最大分页大小，默认 100
- `LOG_LEVEL` — 日志级别 debug/info/warn/error，默认 info
- `API_TOKEN` — API 鉴权 Token，默认 `leaderboard-secret`

## 完整 API 表格

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/boards | 创建榜单 |
| GET | /api/boards | 榜单列表（分页+筛选） |
| GET | /api/boards/{id} | 榜单详情 |
| PUT | /api/boards/{id} | 更新榜单 |
| DELETE | /api/boards/{id} | 删除榜单 |
| POST | /api/members | 创建成员 |
| GET | /api/members | 成员列表（分页+筛选） |
| GET | /api/members/{id} | 成员详情 |
| PUT | /api/members/{id} | 更新成员 |
| DELETE | /api/members/{id} | 删除成员 |
| POST | /api/seasons | 创建赛季 |
| GET | /api/seasons | 赛季列表（分页+筛选） |
| GET | /api/seasons/{id} | 赛季详情 |
| PUT | /api/seasons/{id} | 更新赛季 |
| POST | /api/seasons/{id}/transition | 赛季状态流转 |
| DELETE | /api/seasons/{id} | 删除赛季 |
| POST | /api/score-events | 提交得分事件（自动重算排名） |
| GET | /api/score-events | 得分事件列表（分页+筛选） |
| GET | /api/score-events/{id} | 得分事件详情 |
| PUT | /api/score-events/{id} | 更新得分事件 |
| DELETE | /api/score-events/{id} | 删除得分事件 |
| GET | /api/rank-entries | 排名条目列表（分页+筛选） |
| GET | /api/rank-entries/{id} | 排名条目详情 |
| GET | /api/boards/{board_id}/seasons/{season_id}/top | 榜单 TopN |
| GET | /api/boards/{board_id}/seasons/{season_id}/members/{member_id}/rank | 成员排名 |
| GET | /api/boards/{board_id}/seasons/{season_id}/range | 名次区间查询 |
| DELETE | /api/rank-entries/{id} | 删除排名条目 |
| POST | /api/rank-snapshots | 创建排名快照 |
| GET | /api/rank-snapshots | 快照列表（分页+筛选） |
| GET | /api/rank-snapshots/{id} | 快照详情 |
| DELETE | /api/rank-snapshots/{id} | 删除快照 |
| GET | /api/change-logs | 排名变动历史（分页+筛选） |
| GET | /api/change-logs/{id} | 变动详情 |
| DELETE | /api/change-logs/{id} | 删除变动记录 |
| GET | /api/boards/{board_id}/seasons/{season_id}/export | 导出排名数据（JSON） |

## 实体清单

- Board（榜单）
- Member（成员）
- Season（赛季）
- ScoreEvent（得分事件）
- RankEntry（排名条目）
- RankSnapshot（排名快照）
- ChangeLog（排名变动）
