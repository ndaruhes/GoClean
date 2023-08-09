# STAGE 1
FROM golang:1.19.9-alpine as go-clean-builder
RUN apk add --no-cache pkgconfig
RUN apk add --no-cache gcc musl-dev
RUN apk add --no-cache vips-dev
WORKDIR /app

COPY . /app

WORKDIR ./app
RUN go mod download
CMD cd src && go build -o go-clean

# STAGE 2
FROM alpine:latest
RUN apk add --no-cache gcc musl-dev
RUN apk add --no-cache vips-dev
WORKDIR /root/
COPY --from=go-clean-builder /app .
ENV ENVIRONMENT=local
ENV CONFIG_PATH=/config
EXPOSE 8000
CMD cd src && ./go-clean-builder