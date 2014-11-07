#!/bin/bash -e

cwd=`pwd`

if [[ $1 = "-loc" ]]; then
    find . -name '*.go' | xargs wc -l | sort -n
    exit
fi

VER=0.1.0rc
ID=$(git rev-parse HEAD | cut -c1-7)

cd daemon/actord
go build -ldflags "-X github.com/funkygao/golib/server.VERSION $VER -X github.com/funkygao/golib/server.BuildID $ID -w"
#go build -race -v -ldflags "-X github.com/funkygao/golib/server.VERSION $VER -X github.com/funkygao/golib/server.BuildID $ID -w"

#---------
# show ver
#---------
cd $cwd
./daemon/actord/actord -version
