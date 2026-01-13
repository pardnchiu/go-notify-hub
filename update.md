# Update Log

> Generated: 2026-01-13 

## Recommended Commit Message

```
feat: 實作 Discord webhook 通知 API 服務

建立 Go 專案結構，實作 Discord 通知 API 服務
支援客製化 embed 訊息、多頻道管理與 webhook 配置
```

```
feat: Implement Discord webhook notification API service

Create Go project structure and implement Discord notification API service
Support customizable embed messages, multi-channel management, and webhook configuration
```

***

## Summary

建立 Go 語言通知服務專案，實作 Discord webhook API，支援豐富的 embed 格式，包含標題、描述、圖片、欄位、作者、頁尾等自訂選項，並透過 JSON 檔案管理多個頻道的 webhook URL。

## Changes

### FEAT
- 實作 Discord webhook 通知功能，支援完整的 embed 格式
- 建立 HTTP API 服務，提供 `/discord/:channelName` 端點發送通知
- 實作多頻道管理系統，透過 JSON 檔案配置 webhook URL
- 支援自訂 embed 選項：標題、描述、顏色、時間戳記、圖片、縮圖、欄位、作者、頁尾
- 實作頻道名稱驗證與 webhook URL 快取機制

### CHORE
- 初始化 Go module (`goNotify`) 並配置專案依賴
- 新增 `.gitignore` 忽略系統檔案與 JSON 配置目錄
- 引入 Gin web 框架用於 HTTP 路由處理

***

## Summary

Create a Go notification service project implementing Discord webhook API with rich embed format support, including title, description, images, fields, author, footer, and other customization options, with multi-channel webhook URL management via JSON file.

## Changes

### FEAT
- Implement Discord webhook notification with full embed format support
- Create HTTP API service with `/discord/:channelName` endpoint for sending notifications
- Implement multi-channel management system using JSON file for webhook URL configuration
- Support customizable embed options: title, description, color, timestamp, image, thumbnail, fields, author, footer
- Implement channel name validation and webhook URL caching mechanism

### CHORE
- Initialize Go module (`goNotify`) and configure project dependencies
- Add `.gitignore` to exclude system files and JSON configuration directory
- Integrate Gin web framework for HTTP routing

***

## Files Changed

| File | Status | Tag |
|------|--------|-----|
| `.gitignore` | Added | CHORE |
| `cmd/api/main.go` | Added | FEAT |
| `go.mod` | Added | CHORE |
| `go.sum` | Added | CHORE |
| `internal/bot/line.go` | Added | CHORE |
| `internal/channel/discord.go` | Added | FEAT |
| `internal/handler/discord.go` | Added | FEAT |
