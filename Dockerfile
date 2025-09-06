# build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache git
RUN go run github.com/swaggo/swag/cmd/swag@latest init -g ./cmd/app/main.go -o ./docs -parseDependency -parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/subscriptions ./cmd/app

# final stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/subscriptions /usr/local/bin/subscriptions
COPY --from=builder /app/docs /docs

WORKDIR /
EXPOSE $APP_PORT
ENTRYPOINT ["/usr/local/bin/subscriptions"]
