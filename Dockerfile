FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG APP=api
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/app ./cmd/${APP}

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Jakarta
WORKDIR /app
COPY --from=builder /out/app /app/app
USER nobody
ENTRYPOINT ["/app/app"]
