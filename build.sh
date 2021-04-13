#!/usr/bin/env bash

VERSION=${VERSION:-$(git describe --tags --abbrev=0 2>/dev/null || echo 0.0.0)}
export IMAGE_NAME=${IMAGE_NAME:-"ryts/24coin:${VERSION}"}
export NAMESPACE="test"

BASE_DIR="$(realpath $(dirname "$0"))"
cd $BASE_DIR

OUTPUT_FOLDER="${BASE_DIR}/out"
mkdir -p ${OUTPUT_FOLDER}

while [[ "$#" -gt 0 ]]; do
    case "${1}" in
    --test)
        TEST="Y"
        shift 1
        ;;
    --tview)
        TVIEW="Y"
        shift 1
        ;;
    --build)
        BUILD="Y"
        shift 1
        ;;
    --docker)
        DOCKER="Y"
        shift 1
        ;;
    --kube)
        KUBE="Y"
        shift 1
        ;;
    *)
        echo "Only use '--test', '--tview', '--kube', '--build' and '--docker'" >&2
        exit 1
        ;;
    esac
done

#### Test

if [[ "${TEST}" == "Y" ]]; then
    go test ./... -coverprofile=${OUTPUT_FOLDER}/coverage.out
fi

if [[ "${TVIEW}" == "Y" ]]; then
    go tool cover -html=${OUTPUT_FOLDER}/coverage.out
fi

##########

#### Build

if [[ "${BUILD}" == "Y" ]]; then
    echo "> Build started for ${VERSION}"
    (
        # Build web components
        cd web
        set -e
        pnpm install --prefer-offline
        pnpm run build
        cp -a dist ${BASE_DIR}/internal/api/
        set +e
    )

    echo "> Building app"
    FLAG_V="github.com/rytsh/24coin/internal/common.Version=${VERSION}"

    MAINGO="${BASE_DIR}/cmd/24coin/24coin.go"

    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X ${FLAG_V}" -o ${OUTPUT_FOLDER}/24coin ${MAINGO}
fi

##########

#### Docker Build

if [[ "${DOCKER}" == "Y" ]]; then
    echo "> Building docker"
    tar -czf - out/24coin ci/docker/Dockerfile | docker build -t ${IMAGE_NAME} -f ci/docker/Dockerfile -
fi

##########
