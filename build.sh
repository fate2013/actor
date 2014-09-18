#!/bin/bash -e

cwd=`pwd`

if [[ $1 = "-loc" ]]; then
    find . -name '*.go' | xargs wc -l | sort -n
    exit
fi

cd daemon/actord
ID=$(git rev-parse HEAD | cut -c1-7)
go build -ldflags "-X github.com/funkygao/dragon/server.BuildID $ID -w"
#go build -race -v -ldflags "-X github.com/funkygao/dragon/server.BuildID $ID -w"

#---------
# show ver
#---------
cd $cwd
./daemon/actord/actord -version
