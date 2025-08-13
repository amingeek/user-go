# Builder
FROM golang:1.23-alpine AS builder
WORKDIR /app

# دانلود ماژول‌ها
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct \
 && go mod download

# کپی سورس و ساخت باینری
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /usr/local/bin/user-go main.go

# Runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /usr/local/bin/user-go /usr/local/bin/user-go

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/user-go"]
