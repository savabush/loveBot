# breakfastLoveBot

Telegram bot for a shared breakfast/food menu with carts, stickers, and bilingual UI (EN/RU).

## Features
- Food cards with photos
- Cart flow with accept/partial/decline and approval
- Stickers collection
- Language switcher (EN/RU)

## Requirements
- Go 1.22+
- Telegram Bot token

## Configuration
Set these environment variables (or provide a config file used by `internal/config`):
- `TELEGRAM_BOT_TOKEN`
- `USER_ID1`
- `USER_ID2`

## Run
```bash
go run ./cmd/breakfastLoveBot
```

## Release
Tag a version and push it to GitHub:
```bash
git tag v0.1.0
git push origin v0.1.0
```
GitHub Actions + GoReleaser will build and publish binaries.
