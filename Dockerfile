# --- Stage 1: Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Pertama, hanya copy file go.mod dan go.sum
COPY go.mod go.sum ./

# Ini akan mengambil semua dependensi
# berdasarkan go.mod dan memvalidasinya di dalam lingkungan Linux.
RUN go mod download

# Sekarang, copy sisa kode aplikasi Anda
COPY . .

# Ini akan membuat file go.sum yang 100% valid untuk lingkungan build Linux.
# Langkah ini memastikan tidak ada ketidakcocokan checksum dari mesin Windows Anda.
RUN go mod tidy

# Build aplikasi. Go sekarang akan memiliki semua yang dibutuhkannya.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/web/main.go


# --- Stage 2: Final Image ---
FROM alpine:latest

WORKDIR /app

# Copy binary yang sudah di-build dari stage 'builder'
COPY --from=builder /app/server .

# Copy folder migrasi dari stage 'builder' ke stage final
COPY --from=builder /app/db/migrations ./db/migrations

# Expose port yang digunakan oleh aplikasi Anda di dalam container
EXPOSE 8080

# Command untuk menjalankan aplikasi saat container dimulai
CMD ["./server"]