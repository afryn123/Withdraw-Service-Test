# Withdraw-Service-Test
## Documentation
**1. README.md**

**2. docs/Withdraw-Collection.postman_collection**

## Logs
logs/
*.log

## Database 
Posgresql

## Running Guide

**1. Change file name .env.example to .env and set up**

**2. Install Depedencies**
```bash
go mod tidy 
```
or 
```bash
go mod download
```

**3. Run Migration**
```bash
go run cmd/migrate/main.go
```

**4. Run Apps**
```bash
go run cmd/server/main.go
```
