# Mini Log Engine

Go tabanlı append-only log engine.

## Özellikler
- NDJSON format
- Batch write (10 mesaj)
- Segment rotation (256KB)
- Primary + replica log
- Timeout flush

## Çalıştırma
```bash
go run main.go
