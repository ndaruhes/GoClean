FROM golang:1.19.9-alpine
RUN apk add --no-cache pkgconfig
RUN apk add --no-cache gcc musl-dev
RUN apk add --no-cache vips-dev
WORKDIR /usr/src/app

COPY src .
COPY config ../config
RUN go mod download
RUN go build -o /go-clean

FROM alpine:3.17
RUN apk add --no-cache gcc musl-dev
RUN apk add --no-cache vips-dev
ENV ENVIRONMENT=local
ENV CONFIG_PATH=/config
EXPOSE 8000
COPY --from=go-clean-build /go-clean go-clean

CMD ["./go-clean"]