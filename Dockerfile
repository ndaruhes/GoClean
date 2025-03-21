# STAGE 1
FROM golang:1.23.1-alpine as GoCleanBuilder

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
COPY .env .
COPY go.mod .
COPY go.work .

# Mengunduh dependensi Go
RUN go mod download

# Kompilasi aplikasi Go dan simpan dalam /app/bin/GoClean
RUN go build -o /app/bin/GoClean ./src

# STAGE 2
FROM alpine:latest

# Instalasi dependensi yang diperlukan di tahap kedua
RUN apk add --no-cache gcc musl-dev
RUN apk --no-cache add tzdata
RUN apk add --no-cache vips-dev

# Set env
#ENV APP_ENVIRONMENT=develop
#ENV APP_ROOT_FOLDER=GoClean

# Set direktori kerja di dalam kontainer
WORKDIR /GoClean/

# Menyalin hasil kompilasi dari tahap pertama ke tahap kedua
COPY --from=GoCleanBuilder /app/bin/GoClean ./bin/

# Menyalin konfigurasi yang diperlukan
COPY config ./config/
COPY google ./google/
COPY .env .

# Port yang akan diexpose
EXPOSE 8100

# Perintah untuk menjalankan aplikasi Go
CMD ./bin/GoClean
