# Update Log

> Generated: 2026-01-13 22:23

## Recommended Commit Message

refactor: 重構 Discord handler 程式碼結構，拆分處理函式並提取共用工具
refactor: restructure Discord handler code by splitting handler methods and extracting shared utilities

***

## Summary

重構 Discord webhook 處理器，將原本集中在單一檔案的處理函式拆分為多個獨立檔案，並提取共用的 JSON 寫入功能至 utils 套件。同時新增頻道列表查詢端點，並改善錯誤處理機制。

## Changes

### FEAT
- 新增 `GET /discord/list` 端點，用於列出所有已註冊的 Discord 頻道

### REFACTOR
- 將 `DiscordHandler` 的三個處理方法（Send、Add、Delete）拆分至獨立檔案
  - `internal/handler/discordSend.go` - 處理發送通知
  - `internal/handler/discordAdd.go` - 處理新增頻道
  - `internal/handler/discordDelete.go` - 處理刪除頻道
- 簡化 `discord.go` 主檔案，僅保留結構定義和建構函式
- 提取共用的 `writeJSON` 函式至 `internal/utils/utils.go`，重命名為 `WriteJSON`

### UPDATE
- 修改 `NewDiscordHandler` 建構函式，改為返回錯誤以改善啟動時的錯誤處理
- 更新 `cmd/api/main.go` 以處理建構函式返回的錯誤

***

## Summary

Refactored Discord webhook handler by splitting handler methods from a single file into multiple separate files, and extracted shared JSON writing functionality into a utils package. Added a channel listing endpoint and improved error handling.

## Changes

### FEAT
- Add `GET /discord/list` endpoint to list all registered Discord channels

### REFACTOR
- Split `DiscordHandler` methods into separate files
  - `internal/handler/discordSend.go` - handles sending notifications
  - `internal/handler/discordAdd.go` - handles adding channels
  - `internal/handler/discordDelete.go` - handles deleting channels
- Simplify `discord.go` main file to contain only struct definition and constructor
- Extract shared `writeJSON` function to `internal/utils/utils.go` as `WriteJSON`

### UPDATE
- Modify `NewDiscordHandler` constructor to return error for better startup error handling
- Update `cmd/api/main.go` to handle constructor error

***

## Files Changed

| File | Status | Tag |
|------|--------|-----|
| `cmd/api/main.go` | Modified | FEAT, UPDATE |
| `internal/handler/discord.go` | Modified | REFACTOR |
| `internal/handler/discordSend.go` | Added | REFACTOR |
| `internal/handler/discordAdd.go` | Added | REFACTOR |
| `internal/handler/discordDelete.go` | Added | REFACTOR |
| `internal/utils/utils.go` | Added | REFACTOR |
