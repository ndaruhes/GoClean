# STAGE 1
FROM golang:1.19.9-alpine as go-clean-builder

# Instalasi dependensi yang diperlukan
RUN apk add --no-cache pkgconfig
RUN apk add --no-cache gcc musl-dev
RUN apk add --no-cache vips-dev

# Set direktori kerja di dalam kontainer
WORKDIR /app

# Menyalin seluruh proyek ke dalam kontainer
COPY src ./src
COPY google ./google
COPY config ./config
COPY .air.toml .
COPY go.mod .
COPY go.work .

# Mengunduh dependensi Go
RUN go mod download

# Kompilasi aplikasi Go dan simpan dalam /app/bin/go-clean
RUN go build -o /app/bin/go-clean ./src

# STAGE 2
FROM alpine:latest

# Instalasi dependensi yang diperlukan di tahap kedua
RUN apk add --no-cache gcc musl-dev
RUN apk --no-cache add tzdata
RUN apk add --no-cache vips-dev

# Set direktori kerja di dalam kontainer
WORKDIR /root/

# Menyalin hasil kompilasi dari tahap pertama ke tahap kedua
COPY --from=go-clean-builder /app/bin/go-clean ./

# Menyalin konfigurasi yang diperlukan
COPY config/config-local.yaml ./config/

# Menyalin file Google Config (jika diperlukan)
COPY google/fresh-app.json ./google/

# Port yang akan diexpose
EXPOSE 8000

# Perintah untuk menjalankan aplikasi Go
CMD ./go-clean
