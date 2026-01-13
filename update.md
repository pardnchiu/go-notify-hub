# Update Log

> Generated: 2026-01-13 18:38

## Recommended Commit Message

feat: 新增刪除 Discord 頻道 webhook 的 API 端點
feat: add DELETE endpoint for removing Discord channel webhooks

***

## Summary

新增刪除 Discord 頻道 webhook 的功能，允許透過 DELETE API 端點動態移除已註冊的頻道配置。

## Changes

### FEAT
- 新增 DELETE `/discord/:channelName` API 端點用於刪除 Discord 頻道 webhook
- 實作 `Delete` handler 方法，包含頻道名稱驗證、檔案系統操作及錯誤處理
- 刪除頻道後自動更新 `discord_channel.json` 配置檔

***

## Summary

Added the ability to delete Discord channel webhooks, allowing dynamic removal of registered channel configurations through a DELETE API endpoint.

## Changes

### FEAT
- Add DELETE `/discord/:channelName` API endpoint for removing Discord channel webhooks
- Implement `Delete` handler method with channel name validation, file system operations, and error handling
- Automatically update `discord_channel.json` configuration file after channel deletion

***

## Files Changed

| File | Status | Tag |
|------|--------|-----|
| `cmd/api/main.go` | Modified | FEAT |
| `internal/handler/discord.go` | Modified | FEAT |
