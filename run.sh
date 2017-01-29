#!/bin/bash

#This script is used to run deployment, test of this go app

CMD="$1"
APP="$2"
VERSION="$3"
if [ $CMD == deploy ]; then
    yarn build:prod
    goapp deploy -application "$APP" -version "$VERSION" app.yaml
elif [ $CMD == test ]; then
    go test
fi
