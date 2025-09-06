# build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache git

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/subscriptions ./cmd/app

# final stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/subscriptions /usr/local/bin/subscriptions

WORKDIR /
EXPOSE $APP_PORT
ENTRYPOINT ["/usr/local/bin/subscriptions"]
