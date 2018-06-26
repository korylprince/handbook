#!/bin/bash
set -e

version=$1

tag="korylprince/handbook"

cred=$(cat ~/.git-credentials)

docker build --no-cache --build-arg "CREDENTIALS=$cred" --build-arg "VERSION=$version" --tag "$tag:$version" .

docker push "$tag:$version"

if [ "$2" = "latest" ]; then
    docker tag "$tag:$version" "$tag:latest"
    docker push "$tag:latest"
fi
