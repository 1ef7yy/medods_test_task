FROM golang:alpine AS builder

WORKDIR /build

COPY . .
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/main /app

CMD [ "./main" ]
