![cover](./cover.png)

> [!NOTE]
> 此 README 由 [SKILL](https://github.com/pardnchiu/skill-readme-generate) 生成，英文版請參閱 [這裡](./README.md)。

# go-notify-hub

[![pkg](https://pkg.go.dev/badge/github.com/pardnchiu/go-notify-hub.svg)](https://pkg.go.dev/github.com/pardnchiu/go-notify-hub)
[![card](https://goreportcard.com/badge/github.com/pardnchiu/go-notify-hub)](https://goreportcard.com/report/github.com/pardnchiu/go-notify-hub)
[![license](https://img.shields.io/github/license/pardnchiu/go-notify-hub)](LICENSE)
[![version](https://img.shields.io/github/v/tag/pardnchiu/go-notify-hub?label=release)](https://github.com/pardnchiu/go-notify-hub/releases)

> 多平台通知 API 服務，整合 Discord Webhook、Slack Webhook、LINE Bot 和 Email，透過統一的 RESTful API 管理所有通知管道。

## 目錄

- [功能特點](#功能特點)
- [架構](#架構)
- [安裝](#安裝)
- [設定](#設定)
- [使用方法](#使用方法)
- [API 參考](#api-參考)
- [授權](#授權)
- [Author](#author)
- [Stars](#stars)

## 功能特點

- **多平台整合**：支援 Discord Webhook、Slack Webhook、LINE Bot 和 Email
- **統一 API**：透過 RESTful API 管理所有通知管道
- **管道管理**：動態新增、移除和列出已註冊的管道
- **LINE Bot 互動**：自動處理關注／取消關注事件和批次廣播
- **Email 傳送**：支援單一和批次 Email 傳送，具備 TLS/STARTTLS 支援
- **豐富訊息格式**：支援嵌入內容、附件、欄位、圖片等
- **並發安全**：使用 RWMutex 保護共享資料結構

## 架構

```
cmd/
└── api/
    └── main.go              # 應用程式進入點
internal/
├── bot/
│   ├── dicord/              # Discord Bot 處理邏輯
│   ├── line/                # LINE Bot 處理邏輯
│   └── handler/             # Bot 請求路由
├── channel/
│   ├── discord/             # Discord Webhook 管理
│   └── slack/               # Slack Webhook 管理
├── database/
│   ├── sqlite.go            # SQLite 連線管理
│   ├── insertUser.go        # 新增使用者
│   ├── deleteUser.go        # 刪除使用者
│   └── selectUserLinebot.go # 查詢 LINE Bot 使用者
├── email/
│   ├── email.go             # Email 客戶端初始化
│   ├── send.go              # 傳送單一 Email
│   └── bulk.go              # 批次傳送 Email
└── utils/
    └── utils.go             # 共用工具函式
```

## 安裝

### 從原始碼建置

```bash
# Clone 儲存庫
git clone https://github.com/pardnchiu/go-notify-hub.git
cd go-notify-hub

# 安裝相依套件
go mod download

# 建置應用程式
go build -o go-notify-hub cmd/api/main.go
```

### 使用 Docker

```bash
# 使用 Docker Compose
docker-compose up -d

# 或手動建置 Docker 映像檔
docker build -t go-notify-hub .
docker run -p 8080:8080 \
  -v $(pwd)/data:/data \
  --env-file .env \
  go-notify-hub
```

## 設定

建立 `.env` 檔案於專案根目錄：

```bash
# 資料庫路徑（選填，預設：~/.go-notify-hub/database.db）
DB_PATH=/path/to/database.db

# Discord Bot（選填）
DISCORD_BOT_TOKEN=your_discord_bot_token

# LINE Bot（選填）
LINEBOT_SECRET=your_linebot_channel_secret
LINEBOT_TOKEN=your_linebot_channel_access_token

# Email SMTP（選填）
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_FROM=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
```

### 環境變數說明

| 變數 | 說明 | 必要 | 預設值 |
|------|------|------|--------|
| `DB_PATH` | SQLite 資料庫檔案路徑 | 否 | `~/.go-notify-hub/database.db` |
| `DISCORD_BOT_TOKEN` | Discord Bot Token | 否 | - |
| `LINEBOT_SECRET` | LINE Bot Channel Secret | 否 | - |
| `LINEBOT_TOKEN` | LINE Bot Channel Access Token | 否 | - |
| `EMAIL_SMTP_HOST` | SMTP 伺服器主機 | 否 | - |
| `EMAIL_SMTP_PORT` | SMTP 伺服器埠號 | 否 | `587` |
| `EMAIL_FROM` | 寄件者 Email 地址 | 否 | - |
| `EMAIL_PASSWORD` | Email 密碼或應用程式密碼 | 否 | - |

## 使用方法

### 啟動服務

```bash
# 直接執行
./go-notify-hub

# 或使用 go run
go run cmd/api/main.go
```

服務預設在 `:8080` 啟動。

### Discord Webhook 範例

```bash
# 新增 Discord 管道
curl -X POST http://localhost:8080/discord/add \
  -H "Content-Type: application/json" \
  -d '{
    "name": "alerts",
    "webhook": "https://discord.com/api/webhooks/..."
  }'

# 傳送訊息到 Discord
curl -X POST http://localhost:8080/discord/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello from go-notify-hub!",
    "embeds": [{
      "title": "系統通知",
      "description": "服務已成功啟動",
      "color": 3447003
    }]
  }'

# 列出所有 Discord 管道
curl http://localhost:8080/discord/list

# 刪除 Discord 管道
curl -X DELETE http://localhost:8080/discord/alerts
```

### Slack Webhook 範例

```bash
# 新增 Slack 管道
curl -X POST http://localhost:8080/slack/add \
  -H "Content-Type: application/json" \
  -d '{
    "name": "monitoring",
    "webhook": "https://hooks.slack.com/services/..."
  }'

# 傳送訊息到 Slack
curl -X POST http://localhost:8080/slack/monitoring \
  -H "Content-Type: application/json" \
  -d '{
    "text": "伺服器 CPU 使用率過高",
    "attachments": [{
      "color": "danger",
      "fields": [{
        "title": "CPU 使用率",
        "value": "95%",
        "short": true
      }]
    }]
  }'
```

### LINE Bot 範例

```bash
# Webhook 端點（由 LINE Platform 呼叫）
# POST http://localhost:8080/linebot/webhook

# 廣播訊息給所有關注者
curl -X POST http://localhost:8080/linebot/send/all \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [{
      "type": "text",
      "text": "重要公告：系統將於今晚 22:00 進行維護"
    }]
  }'
```

### Email 範例

```bash
# 傳送單一 Email
curl -X POST http://localhost:8080/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user@example.com"],
    "subject": "測試信件",
    "body": "這是一封測試信件。",
    "html": true
  }'

# 批次傳送 Email
curl -X POST http://localhost:8080/email/send/bulk \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [
      {"email": "user1@example.com", "name": "User 1"},
      {"email": "user2@example.com", "name": "User 2"}
    ],
    "subject": "批次通知",
    "body": "親愛的 {{.Name}}，這是您的個人化訊息。",
    "html": true
  }'
```

## API 參考

### Discord API

| 方法 | 端點 | 說明 |
|--------|----------|-------------|
| `GET` | `/discord/list` | 列出所有已註冊的 Discord 管道 |
| `POST` | `/discord/:channelName` | 傳送訊息到指定管道 |
| `POST` | `/discord/add` | 新增 Discord 管道 |
| `DELETE` | `/discord/:channelName` | 刪除 Discord 管道 |

#### POST /discord/add

**請求 Body：**
```json
{
  "name": "channel-name",
  "webhook": "https://discord.com/api/webhooks/..."
}
```

**回應：**
```json
{
  "message": "Channel added successfully"
}
```

#### POST /discord/:channelName

**請求 Body：**
```json
{
  "content": "Message content",
  "username": "Custom Bot Name",
  "avatar_url": "https://example.com/avatar.png",
  "embeds": [{
    "title": "Embed Title",
    "description": "Embed Description",
    "color": 3447003,
    "fields": [{
      "name": "Field Name",
      "value": "Field Value",
      "inline": true
    }],
    "image": {
      "url": "https://example.com/image.png"
    },
    "timestamp": "2026-01-01T00:00:00Z"
  }]
}
```

### Slack API

| 方法 | 端點 | 說明 |
|--------|----------|-------------|
| `GET` | `/slack/list` | 列出所有已註冊的 Slack 管道 |
| `POST` | `/slack/:channelName` | 傳送訊息到指定管道 |
| `POST` | `/slack/add` | 新增 Slack 管道 |
| `DELETE` | `/slack/:channelName` | 刪除 Slack 管道 |

#### POST /slack/add

**請求 Body：**
```json
{
  "name": "channel-name",
  "webhook": "https://hooks.slack.com/services/..."
}
```

#### POST /slack/:channelName

**請求 Body：**
```json
{
  "text": "Message text",
  "attachments": [{
    "color": "good",
    "title": "Attachment Title",
    "text": "Attachment Text",
    "fields": [{
      "title": "Field Title",
      "value": "Field Value",
      "short": true
    }]
  }]
}
```

### LINE Bot API

| 方法 | 端點 | 說明 |
|--------|----------|-------------|
| `POST` | `/linebot/webhook` | LINE Platform Webhook 端點 |
| `POST` | `/linebot/send/all` | 廣播訊息給所有關注者 |

#### POST /linebot/send/all

**請求 Body：**
```json
{
  "messages": [{
    "type": "text",
    "text": "廣播訊息內容"
  }]
}
```

**支援的訊息類型：**
- `text`：文字訊息
- `image`：圖片訊息
- `video`：影片訊息
- `audio`：音訊訊息
- `location`：位置訊息
- `sticker`：貼圖訊息
- `template`：範本訊息
- `flex`：Flex 訊息

### Email API

| 方法 | 端點 | 說明 |
|--------|----------|-------------|
| `POST` | `/email/send` | 傳送單一 Email |
| `POST` | `/email/send/bulk` | 批次傳送 Email |

#### POST /email/send

**請求 Body：**
```json
{
  "to": ["recipient@example.com"],
  "cc": ["cc@example.com"],
  "bcc": ["bcc@example.com"],
  "subject": "Email Subject",
  "body": "Email body content",
  "html": true
}
```

#### POST /email/send/bulk

**請求 Body：**
```json
{
  "recipients": [
    {"email": "user1@example.com", "name": "User 1"},
    {"email": "user2@example.com", "name": "User 2"}
  ],
  "subject": "Email Subject",
  "body": "Hello {{.Name}}, this is your personalized message.",
  "html": true
}
```

**範本變數：**
- `{{.Name}}`：收件者姓名
- `{{.Email}}`：收件者 Email

### 資料庫 API

#### SQLite

**匯出的類型：**

```go
type SQLite struct {
    // 內部欄位
}
```

**方法：**

| 方法 | 簽章 | 說明 |
|--------|-----------|-------------|
| `NewSQLite` | `func NewSQLite(dbPath string) (*SQLite, error)` | 建立新的 SQLite 連線 |
| `Close` | `func (s *SQLite) Close() error` | 關閉資料庫連線 |
| `InsertUser` | `func (s *SQLite) InsertUser(ctx context.Context, uid string) error` | 新增使用者（LINE Bot 關注者） |
| `DeleteUser` | `func (s *SQLite) DeleteUser(ctx context.Context, uid string) error` | 刪除使用者（LINE Bot 取消關注） |
| `SelectUserLinebot` | `func (s *SQLite) SelectUserLinebot(ctx context.Context) ([]string, error)` | 查詢所有 LINE Bot 使用者 ID |

### 工具函式

**匯出的函式：**

| 函式 | 簽章 | 說明 |
|----------|-----------|-------------|
| `GetPath` | `func GetPath(arg ...string) (string, error)` | 取得設定檔路徑 |
| `GetFile` | `func GetFile(arg ...string) (map[string]string, error)` | 讀取 JSON 設定檔 |
| `WriteJSON` | `func WriteJSON(path string, data map[string]string) error` | 寫入 JSON 設定檔 |
| `ResponseError` | `func ResponseError(c *gin.Context, status int, err error, fn, message string)` | 標準化錯誤回應 |
| `CheckChannelPayload` | `func CheckChannelPayload(req ChannelPayload, regexName, regexWebhook *regexp.Regexp) error` | 驗證管道請求資料 |

## 授權

本專案採用 [MIT LICENSE](LICENSE)。

## Author

<img src="https://avatars.githubusercontent.com/u/25631760" align="left" width="96" height="96" style="margin-right: 0.5rem;">

<h4 style="padding-top: 0">邱敬幃 Pardn Chiu</h4>

<a href="mailto:dev@pardn.io" target="_blank">
<img src="https://pardn.io/image/email.svg" width="48" height="48">
</a> <a href="https://linkedin.com/in/pardnchiu" target="_blank">
<img src="https://pardn.io/image/linkedin.svg" width="48" height="48">
</a>

## Stars

[![Star](https://api.star-history.com/svg?repos=pardnchiu/go-notify-hub&type=Date)](https://www.star-history.com/#pardnchiu/go-notify-hub&Date)

***

©️ 2026 [邱敬幃 Pardn Chiu](https://linkedin.com/in/pardnchiu)
