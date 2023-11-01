#!/usr/bin/env bash

# ---------------------------------------------------------------------
# Globals
# ---------------------------------------------------------------------

DOCKER_REGISTRY="ghcr.io/katasec"
DOCKER_IMAGE_NAME=""
DOCKER_FILE_NAME=""

# ---------------------------------------------------------------------
# Functions
# ---------------------------------------------------------------------
function build() {

    # Set image version from git tag
    IMAGE_VERSION=`git describe --tags --abbrev=0`

    # Check IMAGE_VERSION exists
    if [ -z "${IMAGE_VERSION}" ]; then
        echo "Please specify docker image version via env var IMAGE_VERSION, exitting."
        return 1
    fi

    # Check IMAGE_NAME exists
    if [ -z "${IMAGE_NAME}" ]; then
        echo "Please specify docker image name via env var IMAGE_NAME, exitting."
        return 1
    fi

    # Construct image name
    DOCKER_IMAGE_NAME="${DOCKER_REGISTRY}/${IMAGE_NAME}:${IMAGE_VERSION}"

    # Construct dockerfile name
    DOCKER_FILE_NAME="${IMAGE_NAME}.dockerfile"

    # Build and publish image to artifactory
    cmd="docker build -f ./scripts/${DOCKER_FILE_NAME} -t ${DOCKER_IMAGE_NAME} ."
    echo $cmd
    eval $cmd
    
    if [ $? != 0 ]; then
        return 1
    fi
}

function publishAzure() 
{
    # Rebuild binary into bin folder
    mkdir -p bin
    GOOS=linux GOARCH=amd64 go build -o bin/

    # Publish to Azure function
    func azure functionapp publish katasecweb
}

function runlocal() 
{
    # Rebuild binary into bin folder
    mkdir -p bin
    go build -o bin/

    # Publish to Azure function
    func start
}

function buildAndPush() {
    build
    docker push "${DOCKER_IMAGE_NAME}"
}


# Run the bash function that was passed as a parameter 
$*


