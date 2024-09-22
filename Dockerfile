FROM golang:1.23.0-alpine as builder

RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    && apk add --no-cache bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o pdf-service cmd/app/main.go

FROM alpine:latest

RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    dumb-init

WORKDIR /root/

COPY --from=builder /app/pdf-service .

COPY --from=builder /app/.env .

CMD ["dumb-init", "./pdf-service"]
