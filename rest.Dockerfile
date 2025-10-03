# ---- Tahap 1: Builder ----
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Kompilasi HANYA aplikasi REST
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /rest-server ./cmd/rest

# ---- Tahap 2: Final ----
FROM scratch
COPY --from=builder /rest-server /rest-server

# Port yang diekspos oleh aplikasi REST
EXPOSE 9000

# Perintah untuk menjalankan aplikasi
CMD ["/rest-server"]