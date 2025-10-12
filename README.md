# Golang-Fiber-Tunnel-Wifi-Project
## Documentation
**1. README.md**

**2. docs/Withdraw-Collection.postman_collection**

## Logs
logs/
*.log

## Running Guide

**1. Install Depedencies**
```bash
go mod tidy 
```
or 
```bash
go mod download
```

**2. Run Migration**
```bash
go run cmd/migrate/main.go
```

**4. Chane file name .env.example to .env and set up**

**5. Run Apps**
```bash
go run cmd/server/main.go
```