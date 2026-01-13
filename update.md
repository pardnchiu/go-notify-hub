# Update Log

> Generated: 2026-01-13 18:16

## Recommended Commit Message

feat: 新增動態註冊 Discord 頻道 webhook 的 API 端點
feat: add API endpoint for dynamically registering Discord channel webhooks

***

## Summary

新增 `/discord/add` API 端點以支援動態註冊 Discord 頻道及 webhook URL，並優化頻道配置載入邏輯，確保每次請求都能獲取最新的配置資料。

## Changes

### FEAT
- 新增 `/discord/add` POST 端點，支援批量註冊 Discord 頻道及 webhook URL
- 實作完整的輸入驗證機制，包含頻道名稱格式與 Discord webhook URL 格式檢查
- 支援自動建立 `discord_channel.json` 配置檔案（若不存在）
- 新增批量新增頻道功能，可一次註冊多個頻道與 webhook 對應關係

### FIX
- 修正頻道配置快取問題，改為每次 Send 請求都重新載入 `discord_channel.json`，確保配置資料即時性
- 移除原本的 `channels == nil` 條件判斷，避免配置更新後無法即時生效

### UPDATE
- 新增 Discord webhook URL 驗證正則表達式 `vaildWebhookURL`
- 在 Add 端點中加入 `strings` 套件使用，對輸入資料進行 trim 處理

### REFACTOR
- 抽取 JSON 檔案寫入邏輯至獨立的 `writeJSON` 輔助函式
- 統一檔案寫入格式為 JSON 縮排格式（2 空格）
- 改善錯誤處理與日誌記錄的結構

### STYLE
- 修正 `validChannelName` 正則表達式字串格式，從雙引號改為 backtick

### DOC
- 為 `Send` 函式新增 API 路由註解：`// POST: /discord/send/:channelName`
- 為 `Add` 函式新增 API 路由與請求格式註解

***

## Summary

Added a new `/discord/add` API endpoint to support dynamic registration of Discord channels and webhook URLs, and optimized channel configuration loading logic to ensure the latest configuration data is retrieved on every request.

## Changes

### FEAT
- Add `/discord/add` POST endpoint to support batch registration of Discord channels and webhook URLs
- Implement comprehensive input validation for channel name format and Discord webhook URL format
- Support automatic creation of `discord_channel.json` configuration file if it doesn't exist
- Add batch channel registration functionality to register multiple channel-webhook mappings at once

### FIX
- Fix channel configuration caching issue by reloading `discord_channel.json` on every Send request to ensure real-time configuration data
- Remove original `channels == nil` conditional check to prevent stale configuration after updates

### UPDATE
- Add Discord webhook URL validation regex pattern `vaildWebhookURL`
- Import `strings` package in Add endpoint for trimming input data

### REFACTOR
- Extract JSON file writing logic to standalone `writeJSON` helper function
- Standardize file writing format to indented JSON (2 spaces)
- Improve error handling and logging structure

### STYLE
- Fix `validChannelName` regex string format from double quotes to backticks

### DOC
- Add API route comment for `Send` function: `// POST: /discord/send/:channelName`
- Add API route and request format comments for `Add` function

***

## Files Changed

| File | Status | Tag |
|------|--------|-----|
| `cmd/api/main.go` | Modified | FEAT |
| `internal/handler/discord.go` | Modified | FEAT, FIX, UPDATE, REFACTOR, STYLE, DOC |
