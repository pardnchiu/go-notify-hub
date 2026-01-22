> [!NOTE]
> 此 README 由 [Claude Code](https://github.com/pardnchiu/skill-readme-generate) 生成，英文版請參閱 [這裡](./README.md)。

# go-notification-bot

> 多平台通知機器人服務，整合 Discord、Slack 與 LINE Bot，提供統一的 API 介面發送訊息與管理頻道。

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

- **多平台整合**：支援 Discord Webhook、Slack Webhook 與 LINE Bot
- **統一 API**：透過 RESTful API 管理所有通知頻道
- **Discord Bot 指令**：支援 Slash Command 與傳統訊息指令
- **LINE Bot 互動**：自動處理追蹤/取消追蹤事件與訊息指令
- **股票資訊查詢**：整合 GEX 資料查詢功能
- **批次推播**：LINE Bot 支援最多 500 位使用者的批次訊息推播
- **頻道管理**：動態新增、刪除與列出已註冊頻道

## 架構

```
cmd/
└── api/
    └── main.go          # 程式進入點
internal/
├── channel/
│   ├── discord.go       # Discord Webhook 發送邏輯
│   └── slack.go         # Slack Webhook 發送邏輯
├── database/
│   ├── pg.go            # PostgreSQL 連線管理
│   ├── insertUser.go    # 新增使用者
│   ├── deleteUser.go    # 刪除使用者
│   ├── selectUserLinebot.go  # 查詢 LINE Bot 使用者
│   └── selectTicker.go  # 查詢股票資訊
├── discord/
│   ├── discord.go       # Discord Handler 初始化
│   ├── bot.go           # Discord Bot（Slash Command）
│   ├── send.go          # 發送訊息
│   ├── add.go           # 新增頻道
│   └── delete.go        # 刪除頻道
├── linebot/
│   ├── webhook.go       # LINE Bot Webhook 處理
│   ├── send.go          # 批次推播
│   ├── handleMessage.go # 訊息處理
│   └── commandGex.go    # GEX 指令
├── slack/
│   ├── slack.go         # Slack Handler 初始化
│   ├── send.go          # 發送訊息
│   ├── add.go           # 新增頻道
│   └── delete.go        # 刪除頻道
└── utils/
    └── utils.go         # 通用工具函式
```

## 安裝

```bash
git clone https://github.com/pardnchiu/go-notification-bot.git
cd go-notification-bot
go mod download
```

## 設定

建立 `.env` 檔案並填入以下環境變數：

```env
LINEBOT_SECRET=your_line_channel_secret
LINEBOT_TOKEN=your_line_channel_access_token
DISCORD_TOKEN=your_discord_bot_token
```

## 使用方法

### 啟動伺服器

```bash
go run cmd/api/main.go
```

伺服器將於 `:8080` 啟動。

### Discord Webhook 發送

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
    "description": "伺服器已啟動"
  }'
```

### Slack Webhook 發送

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
    "text": "伺服器已啟動"
  }'
```

### LINE Bot 推播

```bash
# 推播給所有追蹤者
curl -X POST http://localhost:8080/linebot/send/all \
  -H "Content-Type: application/json" \
  -d '{
    "text": "系統公告"
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
  "title": "標題",
  "description": "內容",
  "url": "https://example.com",
  "color": "#FF5733",
  "timestamp": "2025-01-01T00:00:00Z",
  "image": "https://example.com/image.png",
  "thumbnail": "https://example.com/thumb.png",
  "fields": [
    {"name": "欄位名稱", "value": "欄位值", "inline": true}
  ],
  "footer": {"text": "頁尾文字"},
  "author": {"name": "作者名稱"}
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
  "text": "訊息內容",
  "title": "標題",
  "description": "附件內容",
  "color": "#FF5733",
  "image": "https://example.com/image.png",
  "fields": [
    {"title": "欄位名稱", "value": "欄位值", "short": true}
  ],
  "footer": {"text": "頁尾文字"}
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

### Bot 指令

#### Discord

| 指令 | 說明 |
|------|------|
| `/gex <ticker>` | 查詢指定股票的 GEX 資料 |
| `/help` | 顯示可用指令 |

#### LINE Bot

| 指令 | 說明 |
|------|------|
| `/gex $<ticker>` | 查詢指定股票的 GEX 資料 |

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
