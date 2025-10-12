# ---------- Build Stage ----------
FROM golang:1.24-bookworm AS build-stage

ARG APP_DIR=cmd/server

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build binary dari APP_DIR yang dipilih
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -a \
    -o app ./${APP_DIR}/main.go


# ---------- Production Stage ----------
FROM ubuntu:24.04 AS production-stage

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Jakarta

WORKDIR /app

COPY --from=build-stage /app/app .

RUN mkdir -p /app/logs && chmod 777 /app/logs

EXPOSE 8081

CMD ["./app"]