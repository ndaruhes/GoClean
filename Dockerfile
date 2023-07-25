FROM golang:1.19.9-alpine
RUN apk add --no-cache pkgconfig
RUN apk add --no-cache gcc musl-dev
RUN apk add --no-cache vips-dev
WORKDIR /usr/src/app