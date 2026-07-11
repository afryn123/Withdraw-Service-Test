# Withdraw Service 

## Struktur

```text
cmd/
  api/                       composition root dan HTTP server
  migrate/                   database migration
internal/
  auth/
  transactions/
  user/
  wallet/
    application/             use case/orchestration
    delivery/                adapter HTTP Gin
    domain/                  entity, business rule, dan port/interface
    infrastructure/          implementasi GORM/crypto/generator
  shared/                    config, transaction, middleware, response
  mocks/                     hasil generate GoMock
```

Dependency mengarah ke domain/application. Application service tidak bergantung
pada Gin atau GORM; dependency konkret dirangkai hanya di `cmd/api/main.go`.

## Menjalankan

Salin `.env.example` menjadi `.env`, lalu:

```bash
go mod download
go run ./cmd/migrate
go run ./cmd/api
```

Atau jalankan seluruh stack:

```bash
docker compose up --build
```

API menangani `SIGINT` dan `SIGTERM` dengan graceful shutdown selama maksimal
10 detik sebelum koneksi HTTP dan database ditutup.

## Logging

Aplikasi menggunakan `log/slog` dengan format JSON dan menulis ke stdout. Setiap
request memperoleh `X-Request-ID` dan menghasilkan field `method`, `path`,
`status`, `latency_ms`, dan `client_ip`. Level dapat diatur melalui `LOG_LEVEL`
(`debug`, `info`, `warn`, atau `error`). Aplikasi tidak membuat atau merotasi
file log.

Docker Compose memakai driver `json-file` dengan rotasi maksimal 10 MB per file
dan lima file. Pada Kubernetes, stdout/stderr dibaca oleh container runtime;
retensi dan pengiriman log sebaiknya dikonfigurasi pada node atau collector
seperti Fluent Bit, Vector, Promtail/Loki, atau platform logging yang digunakan.

## OpenAPI dan Swagger UI

Setelah API berjalan, dokumentasi tersedia di:

- `GET /docs/` untuk Swagger UI.
- `GET /openapi.yaml` untuk spesifikasi OpenAPI 3.0.3.

OpenAPI menggunakan server URL relatif (`.`), sehingga tombol **Try it out**
otomatis mengikuti scheme, host, port, dan path prefix tempat dokumentasi dibuka.
Contohnya, dokumentasi pada `https://example.com/withdraw/docs/` akan memanggil
API di bawah `https://example.com/withdraw/`.

File OpenAPI di-embed ke binary. Tampilan Swagger UI memuat asset versi 5 dari
CDN `unpkg.com`, sehingga browser membutuhkan akses internet untuk menampilkan UI;
spesifikasi `/openapi.yaml` tetap tersedia tanpa akses CDN.

## Test

Unit test menggunakan Testify untuk assertion dan GoMock untuk dependency mock.

```bash
go generate ./internal/mocks
go test ./...
go test -cover ./internal/... 
```

Integration test repository menggunakan PostgreSQL terpisah dan hanya berjalan
bila build tag serta DSN test diberikan:

```bash
TEST_DATABASE_DSN="host=127.0.0.1 user=postgres password=password dbname=withdraw_test port=5432 sslmode=disable" \
go test -tags=integration ./internal/user/infrastructure
```

Endpoint tetap mengikuti service awal:

- `POST /api/auth/login`
- `POST /api/users/create`
- `PATCH /api/users/:user_id` (Bearer token, hanya profil sendiri)
- `DELETE /api/users/:user_id` (Bearer token, hanya akun sendiri)
- `GET /api/wallet/:user_id/balance`
- `POST /api/transaction/withdraw`
- `GET /health`
- `GET /docs/`
- `GET /openapi.yaml`
