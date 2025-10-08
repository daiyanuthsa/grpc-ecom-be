# ---- Tahap 1: Builder ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Salin file modul untuk men-download dependensi
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode sumber
COPY . .

# Kompilasi HANYA aplikasi gRPC
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /grpc-server ./cmd/grpc

# ---- Tahap 2: Final ----
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y ca-certificates

# Salin HANYA file executable yang sudah jadi
COPY --from=builder /grpc-server /grpc-server

# Port yang diekspos oleh aplikasi gRPC
EXPOSE 50051

# Perintah untuk menjalankan aplikasi
CMD ["/grpc-server"]