#!/bin/bash

version=$1

cred=$(cat ~/.git-credentials)

docker build --no-cache --build-arg "CREDENTIALS=$cred" --build-arg "VERSION=$version" --tag "korylprince/handbook:$version" .

docker push "korylprince/handbook:$version"
