# ---------- BUILD ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o calendar-wallpaper

# ---------- RUNTIME ----------
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/calendar-wallpaper .
COPY --from=builder /app/fonts ./fonts
COPY --from=builder /app/web ./web

ENV PORT=8080

EXPOSE 8080

CMD ["./calendar-wallpaper"]
