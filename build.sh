#!/bin/bash -e

if [[ $1 = "-loc" ]]; then
    find . -name '*.go' | xargs wc -l | sort -n
    exit
fi

ID=$(git rev-parse HEAD | cut -c1-7)
go build -ldflags "-X github.com/funkygao/dragon/server.BuildID $ID -w"

#---------
# show ver
#---------
./daemon/dragon/dragon -version
