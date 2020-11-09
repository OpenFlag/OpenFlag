#!/usr/bin/env bash

APP=$1
APP_VERSION=$2
LDFLAGS=$3

RELEASE_FOLDER=releases

function release {
    GO_OS=$1
    GO_ARCH=$2

    GOOS=${GO_OS} GOARCH=${GO_ARCH} go build -ldflags "${LDFLAGS}"  ./cmd/${APP}

    FOLDER_NAME=${APP}-${APP_VERSION}-${GO_OS}-${GO_ARCH}
    mkdir ${FOLDER_NAME}

    mv ${APP} ${FOLDER_NAME}
    cp LICENSE ${FOLDER_NAME}
    cp README.md ${FOLDER_NAME}
    cp config.example.yaml ${FOLDER_NAME}/config.yaml

    mkdir -p ${FOLDER_NAME}/browser/${APP}-ui
    cp -r ./browser/${APP}-ui/build ${FOLDER_NAME}/browser/${APP}-ui

    zip -r ${FOLDER_NAME}.zip ${FOLDER_NAME}
    rm -rf ${FOLDER_NAME}
    mv ${FOLDER_NAME}.zip ${RELEASE_FOLDER}
}

mkdir -p ${RELEASE_FOLDER}

cd ./browser/${APP}-ui
npm install
npm run build --prod
cd ../../

release darwin amd64
release linux amd64
