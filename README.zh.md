![cover](./cover.png)

> [!NOTE]
> 此 README 由 [Claude Code](https://github.com/pardnchiu/skill-readme-generate) 生成，英文版請參閱 [這裡](./README.md)。

# go-notify-hub

> 多平台通知 API 服務，整合 Discord Webhook、Slack Webhook、LINE Bot 與 Email 發送功能，透過統一的 RESTful API 管理所有通知頻道。

## 目錄

- [功能特點](#功能特點)
- [架構](#架構)
- [安裝](#安裝)
- [設定](#設定)
- [使用方法](#使用方法)
- [API 參考](#api-參考)
- [授權](#授權)
- [Author](#author)

## 功能特點

- **多平台整合**：支援 Discord Webhook、Slack Webhook、LINE Bot 與 Email
- **統一 API**：透過 RESTful API 管理所有通知頻道
- **頻道管理**：動態新增、刪除與列出已註冊頻道
- **LINE Bot 互動**：自動處理追蹤/取消追蹤事件與批次推播
- **Email 發送**：支援單封與批量郵件發送，含 TLS/STARTTLS
- **豐富訊息格式**：支援 Embeds、附件、欄位、圖片等進階格式
- **併發安全**：使用 RWMutex 保護共享資料結構

## 架構

```
cmd/
└── api/
    └── main.go              # 程式進入點
internal/
├── channel/
│   ├── discord.go           # Discord Webhook 發送邏輯
│   └── slack.go             # Slack Webhook 發送邏輯
├── database/
│   ├── sqlite.go            # SQLite 連線管理
│   ├── insertUser.go        # 新增使用者
│   ├── deleteUser.go        # 刪除使用者
│   └── selectUserLinebot.go # 查詢 LINE Bot 使用者
├── discord/
│   ├── discord.go           # Discord Handler 初始化
│   ├── send.go              # 發送訊息
│   ├── add.go               # 新增頻道
│   └── delete.go            # 刪除頻道
├── email/
│   ├── email.go             # Email 客戶端與 SMTP 發送
│   ├── send.go              # 單封郵件發送
│   └── bulk.go              # 批量郵件發送
├── linebot/
│   ├── webhook.go           # LINE Bot Webhook 處理
│   ├── send.go              # 批次推播
│   └── handleMessage.go     # 訊息處理
├── slack/
│   ├── slack.go             # Slack Handler 初始化
│   ├── send.go              # 發送訊息
│   ├── add.go               # 新增頻道
│   └── delete.go            # 刪除頻道
└── utils/
    └── utils.go             # 通用工具函式
```

## 安裝

### 前置需求

- Go 1.20 或更高版本

### 下載與編譯

```bash
git clone https://github.com/pardnchiu/go-notify-hub.git
cd go-notify-hub
go mod download
go build -o go-notify-hub ./cmd/api
```

## 設定

建立 `.env` 檔案並填入以下環境變數：

```env
# 資料庫路徑（選填，預設為 ~/.go-notify-hub/database.db）
DB_PATH=/path/to/database.db

# LINE Bot（選填）
LINEBOT_SECRET=your_line_channel_secret
LINEBOT_TOKEN=your_line_channel_access_token

# Email SMTP（必填，若使用 Email 功能）
MAIL_SERVICE=smtp.example.com
MAIL_SERVICE_PORT=587
MAIL_SERVICE_USER=user@example.com
MAIL_SERVICE_PASSWORD=your_password
```

## 使用方法

### 啟動伺服器

```bash
go run cmd/api/main.go
```

伺服器將於 `:8080` 啟動。

### Docker

```bash
docker-compose up -d
```

### Discord Webhook

```bash
# 新增頻道
curl -X POST http://localhost:8080/discord/add \
  -H "Content-Type: application/json" \
  -d '{
    "datas": [
      {"name": "alerts", "webhook": "https://discord.com/api/webhooks/..."}
    ]
  }'

# 發送訊息
curl -X POST http://localhost:8080/discord/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "系統通知",
    "description": "伺服器已啟動",
    "color": "#00FF00"
  }'

# 列出頻道
curl http://localhost:8080/discord/list

# 刪除頻道
curl -X DELETE http://localhost:8080/discord/alerts
```

### Slack Webhook

```bash
# 新增頻道
curl -X POST http://localhost:8080/slack/add \
  -H "Content-Type: application/json" \
  -d '{
    "datas": [
      {"name": "general", "webhook": "https://hooks.slack.com/services/..."}
    ]
  }'

# 發送訊息
curl -X POST http://localhost:8080/slack/general \
  -H "Content-Type: application/json" \
  -d '{
    "text": "系統通知",
    "title": "伺服器狀態",
    "description": "所有服務運作正常",
    "color": "good"
  }'
```

### LINE Bot 推播

```bash
# 推播給所有追蹤者
curl -X POST http://localhost:8080/linebot/send/all \
  -H "Content-Type: application/json" \
  -d '{
    "text": "系統公告",
    "image": "https://example.com/image.png"
  }'
```

### Email 發送

```bash
# 單封郵件
curl -X POST http://localhost:8080/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": "recipient@example.com",
    "subject": "測試郵件",
    "body": "這是一封測試郵件",
    "is_html": false
  }'

# 批量郵件
curl -X POST http://localhost:8080/email/send/bulk \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user1@example.com", "user2@example.com"],
    "subject": "系統公告",
    "body": "<h1>公告內容</h1>",
    "is_html": true,
    "min_delay": 1
  }'
```

## API 參考

### Discord

| 方法 | 路徑 | 說明 |
|------|------|------|
| GET | `/discord/list` | 列出所有已註冊頻道 |
| POST | `/discord/add` | 新增頻道 |
| POST | `/discord/:channelName` | 發送訊息至指定頻道 |
| DELETE | `/discord/:channelName` | 刪除頻道 |

#### Discord 訊息格式

```json
{
  "title": "標題（必填）",
  "description": "內容（必填）",
  "url": "https://example.com",
  "color": "#FF5733",
  "timestamp": "2025-01-01T00:00:00Z",
  "image": "https://example.com/image.png",
  "thumbnail": "https://example.com/thumb.png",
  "fields": [
    {"name": "欄位名稱", "value": "欄位值", "inline": true}
  ],
  "footer": {"text": "頁尾文字", "icon_url": "https://example.com/icon.png"},
  "author": {"name": "作者名稱", "url": "https://example.com", "icon_url": "https://example.com/author.png"},
  "username": "自訂 Bot 名稱",
  "avatar_url": "https://example.com/avatar.png"
}
```

### Slack

| 方法 | 路徑 | 說明 |
|------|------|------|
| GET | `/slack/list` | 列出所有已註冊頻道 |
| POST | `/slack/add` | 新增頻道 |
| POST | `/slack/:channelName` | 發送訊息至指定頻道 |
| DELETE | `/slack/:channelName` | 刪除頻道 |

#### Slack 訊息格式

```json
{
  "text": "訊息文字（必填）",
  "title": "Attachment 標題",
  "title_link": "https://example.com",
  "description": "Attachment 內容",
  "pretext": "Attachment 上方文字",
  "color": "#FF5733",
  "timestamp": 1704067200,
  "image": "https://example.com/image.png",
  "thumbnail": "https://example.com/thumb.png",
  "fields": [
    {"title": "欄位名稱", "value": "欄位值", "short": true}
  ],
  "footer": {"text": "頁尾文字", "icon_url": "https://example.com/icon.png"},
  "username": "自訂 Bot 名稱",
  "icon_emoji": ":rocket:",
  "channel": "#channel",
  "thread_ts": "1234567890.123456"
}
```

### LINE Bot

| 方法 | 路徑 | 說明 |
|------|------|------|
| POST | `/linebot/webhook` | LINE Webhook 端點 |
| POST | `/linebot/send/all` | 推播給所有追蹤者 |

#### LINE Bot 訊息格式

```json
{
  "text": "訊息內容",
  "image": "https://example.com/image.png",
  "image_preview": "https://example.com/preview.png"
}
```

### Email

| 方法 | 路徑 | 說明 |
|------|------|------|
| POST | `/email/send` | 發送單封郵件 |
| POST | `/email/send/bulk` | 批量發送郵件 |

#### Email 發送格式

```json
{
  "to": "recipient@example.com",
  "subject": "郵件主旨（必填）",
  "body": "郵件內容（必填）",
  "alt_body": "純文字替代內容",
  "from": "sender@example.com",
  "cc": "cc@example.com",
  "bcc": "bcc@example.com",
  "priority": "high",
  "is_html": true
}
```

#### Email 批量發送格式

```json
{
  "to": ["user1@example.com", "user2@example.com"],
  "subject": "郵件主旨（必填）",
  "body": "郵件內容（必填）",
  "from": "sender@example.com",
  "is_html": false,
  "min_delay": 1,
  "stop_on_error": false
}
```

**收件人格式支援**：

- 單一地址：`"user@example.com"`
- 多個地址：`["user1@example.com", "user2@example.com"]`
- 帶名稱：`"Name:user@example.com"` 或 `{"user@example.com": "Name"}`

## 授權

MIT License

## Author

<img src="https://avatars.githubusercontent.com/u/25631760" align="left" width="96" height="96" style="margin-right: 0.5rem;">

<h4 style="padding-top: 0">邱敬幃 Pardn Chiu</h4>

<a href="mailto:dev@pardn.io" target="_blank">
<img src="https://pardn.io/image/email.svg" width="48" height="48">
</a> <a href="https://linkedin.com/in/pardnchiu" target="_blank">
<img src="https://pardn.io/image/linkedin.svg" width="48" height="48">
</a>

***

©️ 2026 [邱敬幃 Pardn Chiu](https://linkedin.com/in/pardnchiu)
