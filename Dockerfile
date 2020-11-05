#
# 1. OpenFlag Build Container
#
FROM golang:1.15-alpine AS openflag-backend

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

ARG GO_PROXY
ENV GOPROXY=${GO_PROXY}

RUN mkdir -p /src

RUN apk add build-base git

COPY go.sum go.mod /src/

WORKDIR /src

COPY . /src

RUN make install && \
    mkdir -p /app && \
    cp -r $GOPATH/bin/openflag /app/

#
# 2. OpenFlag UI Container
#
FROM node:alpine AS openflag-ui

WORKDIR /usr/src/browser/openflag-ui

COPY browser/openflag-ui/package.json browser/openflag-ui/package-lock.json ./

RUN npm install

ADD /browser/openflag-ui ./

RUN npm run build --prod

#
# 3. Runtime Container
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

COPY --from=openflag-backend /app /app/

COPY --from=openflag-backend /src/docker/run.sh /app/run.sh

COPY --from=openflag-backend /src/internal/app/openflag/migrations /app/migrations

COPY --from=openflag-ui /usr/src/openflag-ui /app/browser/openflag-ui

CMD ['run.sh']
