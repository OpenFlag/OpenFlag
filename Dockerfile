#
# 1. Build Container
#
FROM golang:1.15-alpine AS build

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

ARG GO_PROXY
ENV GOPROXY=${GO_PROXY}

RUN mkdir -p /src

RUN apk add build-base
RUN apk add git

# First add modules list to better utilize caching
COPY go.sum go.mod /src/

WORKDIR /src

COPY . /src

# Build components.
# Put built binaries and runtime resources in /app dir ready to be copied over or used.
RUN make install && \
    mkdir -p /app && \
    cp -r $GOPATH/bin/gopher /app/

#
# 2. Runtime Container
#
FROM alpine:3.9

ENV TZ=Asia/Tehran \
    PATH="/app:${PATH}"

RUN apk add --update --no-cache \
      tzdata \
      ca-certificates \
      bash \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    mkdir -p /var/log && \
    chgrp -R 0 /var/log && \
    chmod -R g=u /var/log

WORKDIR /app

COPY --from=build /app /app/

COPY --from=build /src/migrations /app/migrations
