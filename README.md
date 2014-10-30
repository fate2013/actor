actor
=====

                _             
               | |            
      __ _  ___| |_ ___  _ __ 
     / _` |/ __| __/ _ \| '__|
    | (_| | (__| || (_) | |   
     \__,_|\___|\__\___/|_|   
    
### Install

    # install golang amd64 
    https://golang.org/dl/

    # setup GOPATH env
    export GOPATH=~/gopkg

    # install actor package
    go get github.com/funkygao/dragon
    cd $GOPATH/src/github.com/funkygao/dragon

    # install all dependencies
    go get ./... 

    # build the executable
    ./build.sh

    # create a config file
    cp etc/actord.cf.sample etc/actord.cf (change the address to your own)

    # startup the daemon
    ./daemon/actord/actord

### actor IS

* external scheduler for delayed jobs
  - PvP march
  - PvE march
  - Job

* serializer for concurrent updates
  - lock maintainer and issuer with retry mechanism
  - actor make concurrent calls into sequential calls

* coodinator
  - everything that may lead to race condition and concurrent updates will be decided by actor

### TODO
* March may need K
* alliance lock
* pprof may influnce performance
* mysql transaction with isolation repeatable read + optimistic locking has same effect
  - I'd rather kill actor instead of mysqld
  - what about distributed mysql instances?
* teleport
*   Write/Read timeout and check N in loop
*   can a player send N marches to the same tile?
*   simulate mysql shutdown
    - done! golang mysql driver with breaker will handle this
*   WHERE UNIX_TIMESTAMP(time_end) index hit
    - need to optimize DB index
*   worker throttle
    - we can't have toooo many callbacks concurrently, use channel for throttle, easy...
*   handles NULL column
    - march.type done if it's NULL, what about others?
    - maybe we should let DB handle this
    - but mysql enum datatype can't handle this automatically
*   tsung 20M rows in db, and try actor

