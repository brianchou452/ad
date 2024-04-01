#!/bin/bash -ex

CMD=${1}

IMG_NAME=dcard-ad-backend
BUILD_IMG_NAME=ghcr.io/brianchou452/dcard-ad-backend
DOCKER_COMPOSE_NAMESPACE=$IMG_NAME

if [ "$CMD" == "BUILD" -a "$#" == "2" ]; then
    BUILD_VERSION=${2}
    echo "Run Build(BUILD_VERSION=$BUILD_VERSION)"
    docker build --pull\
                 --build-arg BUILD_VERSION=$BUILD_VERSION\
                 -t $IMG_NAME\
                 .

elif [ "$CMD" == "TO_GITHUB" -a "$#" == "2" ]; then
    BUILD_VERSION=${2}
    echo "Run Build(BUILD_VERSION=$BUILD_VERSION)"
    docker buildx build --pull\
                        --build-arg BUILD_VERSION=$BUILD_VERSION\
                        --platform=linux/amd64\
                        --push\
                        -t $BUILD_IMG_NAME:$BUILD_VERSION\
                        .

elif [ "$CMD" == "DEV_UP" -a  $# -le 2 ]; then
    REBUILD=${2:-ok}

    echo "Run DEV_UP(REBUILD=$REBUILD)"

    if [ "$REBUILD" == "ok" ]; then
        ./scripts/cicd.sh BUILD dev
    fi

    docker compose -p $DOCKER_COMPOSE_NAMESPACE\
                   -f docker/compose-dev.yaml\
                   up -d

elif [ "$CMD" == "DEV_DOWN"  -a  "$#" == "1" ]; then
    echo "Run DEV_DOWN()"
    docker compose -p $DOCKER_COMPOSE_NAMESPACE\
                   -f docker/compose-dev.yaml\
                   down --remove-orphans

elif [ "$CMD" == "PRESSURE_TEST"  -a  "$#" == "1" ]; then
    echo "Run PRESSURE_TEST()"
    docker compose -p $DOCKER_COMPOSE_NAMESPACE\
                   -f docker/compose-dev.yaml\
                   up -d
    
    docker compose -p $DOCKER_COMPOSE_NAMESPACE\
                   -f docker/compose-dev.yaml\
                   -f docker/compose-pressure-test.yaml\
                   run k6 run /scripts/main.js

else
    echo "Invalid command"
    exit 1
fi
