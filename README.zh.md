> [!NOTE]
> 此 README 由 Claude Code 生成，英文版請參閱 [這裡](./README.md)。

# go-notification-bot

[![Go](https://img.shields.io/badge/Go-1.25.1-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Gin](https://img.shields.io/badge/Gin-1.11.0-00ADD8)](https://gin-gonic.com)

> 多頻道通知服務，整合 Discord、Slack 與 LINE Bot，提供統一的 REST API 介面進行訊息推送與 Webhook 管理。

## 功能特點

- **多平台支援**：整合 Discord Webhook、Slack Webhook 與 LINE Bot Messaging API
- **統一 API**：透過 REST API 管理所有通知頻道
- **Webhook 管理**：支援新增、刪除、列表查詢頻道設定
- **LINE Bot 指令**：支援自訂指令解析（如 `/gex $TICKER`）
- **PostgreSQL 整合**：使用者管理與資料持久化

## 架構

```mermaid
graph TB
    subgraph Client
        A[HTTP Request]
    end

    subgraph API["Gin Server :8080"]
        B[Router]
        B --> C[Discord Handler]
        B --> D[Slack Handler]
        B --> E[LINE Bot Handler]
    end

    subgraph Storage
        F[(PostgreSQL)]
        G[JSON Files]
    end

    subgraph External["外部服務"]
        H[Discord Webhook]
        I[Slack Webhook]
        J[LINE Messaging API]
    end

    A --> B
    C --> G
    C --> H
    D --> G
    D --> I
    E --> F
    E --> J
```

### 目錄結構

```
go-notification-bot/
├── cmd/api/           # 應用程式進入點
├── internal/
│   ├── channel/       # 訊息發送邏輯（Discord、Slack）
│   ├── database/      # PostgreSQL 資料庫操作
│   ├── discord/       # Discord handler
│   ├── linebot/       # LINE Bot handler
│   ├── slack/         # Slack handler
│   └── utils/         # 共用工具函式
└── json/              # 頻道設定檔儲存
```

### 請求流程

```mermaid
sequenceDiagram
    participant C as Client
    participant G as Gin Router
    participant H as Handler
    participant S as Storage
    participant E as External Service

    C->>G: HTTP Request
    G->>H: Route to Handler
    H->>S: Read/Write Config
    H->>E: Send Notification
    E-->>H: Response
    H-->>C: JSON Response
```

## 安裝

```bash
git clone https://github.com/pardnchiu/go-notification-bot.git
cd bot
go mod download
```

## 設定

建立 `.env` 檔案：

```env
LINEBOT_SECRET=your_line_bot_secret
LINEBOT_TOKEN=your_line_bot_token
DATABASE_URL=postgres://user:password@localhost:5432/dbname
```

## 執行

```bash
go run cmd/api/main.go
```

或使用 Docker：

```bash
docker-compose up
```

## API 端點

### Discord

| 方法 | 路徑 | 說明 |
|------|------|------|
| GET | `/discord/list` | 列出所有已註冊的 Discord 頻道 |
| POST | `/discord/add` | 新增 Discord Webhook 頻道 |
| DELETE | `/discord/delete/:channelName` | 刪除指定頻道 |
| POST | `/discord/send/:channelName` | 發送訊息至指定頻道 |

### Slack

| 方法 | 路徑 | 說明 |
|------|------|------|
| GET | `/slack/list` | 列出所有已註冊的 Slack 頻道 |
| POST | `/slack/add` | 新增 Slack Webhook 頻道 |
| DELETE | `/slack/delete/:channelName` | 刪除指定頻道 |
| POST | `/slack/send/:channelName` | 發送訊息至指定頻道 |

### LINE Bot

| 方法 | 路徑 | 說明 |
|------|------|------|
| POST | `/linebot/webhook` | LINE Bot Webhook 接收端點 |
| POST | `/linebot/send/all` | 群發訊息給所有訂閱用戶 |

## 使用範例

### 新增 Discord 頻道

```bash
curl -X POST http://localhost:8080/discord/add \
  -H "Content-Type: application/json" \
  -d '{
    "datas": [{
      "name": "general",
      "webhook": "https://discord.com/api/webhooks/..."
    }]
  }'
```

### 發送 Discord 訊息

```bash
curl -X POST http://localhost:8080/discord/send/general \
  -H "Content-Type: application/json" \
  -d '{
    "title": "通知標題",
    "description": "訊息內容"
  }'
```

### LINE Bot 群發訊息

```bash
curl -X POST http://localhost:8080/linebot/send/all \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, World!",
    "image": "https://example.com/image.png",
    "image_preview": "https://example.com/image_preview.png"
  }'
```

## 授權

此專案為私人專案。

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
