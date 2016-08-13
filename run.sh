#!/bin/bash

#This script is used to run deployment, test of this go app

CMD="$1"
if [ $CMD == deploy ]; then
    $HOME/go_appengine/appcfg.py update "$PWD"
elif [ $CMD == test ]; then
    go test
fi
