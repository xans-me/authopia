# syntax = docker/dockerfile:1.0-experimental
FROM golang:alpine AS builder

LABEL maintainer="mulia.ichsan17@gmail.com"

ENV GO111MODULE=on
ENV GOBIN /go/bin
ENV GOPATH /app
ENV env ${APP_ENV}

RUN mkdir -p /app/authopia
WORKDIR /app/authopia
ADD . /app/authopia


# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM alpine:3.11

ARG APP_ENV
#ENV env ${APP_ENV}
WORKDIR /app/authopia

RUN apk update && apk add --no-cache tzdata ca-certificates curl && \
    cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    echo "Asia/Jakarta" > /etc/timezone

COPY --from=builder /app/authopia/main /app/authopia/main
COPY --from=builder /app/authopia/env /app/authopia/env
COPY --from=builder /app/authopia/core /app/authopia/core


EXPOSE 3000
#CMD ./main -env $env
CMD ./main
