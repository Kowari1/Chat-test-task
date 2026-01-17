FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/server ./server

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
