# ---------- STAGE 1: build ----------
FROM golang:1.25-alpine AS builder

WORKDIR /src/app

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app/ .

# Сборка статического бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio main.go

# ---------- STAGE 2: run ----------
FROM scratch

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /src/app/portfolio .

# Копируем статические файлы
COPY app/static ./static
COPY app/data.yml ./
COPY app/templates ./templates

EXPOSE 8080
CMD ["./portfolio"]
