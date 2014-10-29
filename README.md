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

### TODO

* Job may race with March/PVEMarch
* Delayed Job race with player instant actions
* 


### Introduction

RTS: synchronization is king

RTS game non-instant scheduler to handle the following challanges:
* concurrency isolation
* scheduling of delayed jobs

Non-instant events inclues:
* Job
  - quest
  - construction
  - heal troops
  - train troop
  - troop boost
  - vip
  - upkeep reduce
  - research
  - resource boost
* March
  - each march status, e,g. arrive, goHome, gathering
* PVEMarch

### Assumptions
*   concurrency(time) accuracy is in 1s
    - there is no absolute order although in real world we can't step into the same rive at the same time
    - it means if 2-more events happens at the same second, they are concurrent
    - concurrent events will be serialized by actor or by chance(lock, depends on who accquired the lock first)
*   combat is atomic and always between 2 parties(always A@Z instead of A,B,C@Z)
    - combat must lock attacker & attackee
*   actor simplifies concurrency problem into 1-2 player problem
    - actor side locking
    - actor make concurrent calls into sequential calls
    - actor is the coordinator
*   we can use 'divide and conquer' methology to simplify problems

### Benefits of introducing actor
*   actor make N-N problem -> N-1 -> 2-1 player problem
    - so that in php, we can lock only attacker/attackee uid
*   the only instant multiple user action problem is: teleport
    - it has problem because of (instant/job conflicts)
    - with help of DB integrity, it gets solved
    - imagine the scenario

### Lock, lock what?
*   just user lock: current user(instant action)
*   attacker and attachee
*   any action may fail because of the mutex lock
    - job/march fail: actor will retry after 1s
    - instant action fail: client will popup dialog to inform user: user decide when to retry

### Problems
*   server push lost and out-of-order
    - war event
    - alliance event
*   race condition between instant action and job based action

A gathering, 2s left, then B attack A, push is late,so A loadMarch to server, but at this
moment, A might go home or arrive home. Then the push arrives, how client handle this push?

(precondition, effect) - atomic pair

teleport race with actor

pve march callback, then the march state changed, lock fails, -> 2 concurrent callbacks for the 1 march

Job and March for a player wakes at the same time, should be merged into one callback
e,g HealTroop job and opponent march arrives at the same time?

race condition will have effects on:
1. combat result : hero equip, heal troops
2. pillagings    : resource

    
                    system
                      |
                +--  tile    --+
                |              |
                |              |
    (producer)  +-- units    --+  (consumer)
       player --|-- resource --|-- opponent
                +-- hero     --+
                      |
                      |
                    actor
    


                   swoole    +-- player(rw)
              +--- php-fpm --|
              |              +-- actor (rw)
              |
        db ---|--- mq worker(ro)
              |
              |
              +--- batch script(resource tile generator)


        
                            +- Job
                            |- CityMap
                  +- Job ---|- Encounter
                  |         |- AllianceVipQuest
                  |         +- Research (why use this table?)
                  |
        Wakeable -|- March
                  |
                  +- PVEMarch
        
        
                        +- speedup  ----+               +- consume
                        |- boost    ----|               |- build
                        |- recall   ----|               |- 
    march -+    actor   |- combat   ----| race with     |-
            |-----------|               |---------------|-
    job   -+            |               |               |
                        +-             -+               +



critical setions solved by actor:
1. tile enamping conflicts(concurrently or neighbore)
2. multiple opponent attack the same city
3. before do something, we no longer need to wake up player's job(e,g. before combat, wakeup attacker/attackee's job?)

### err handling
[MySQL] 2014/09/24 09:29:15 packets.go:30: read tcp 192.168.42.106:3306: operation timed out
[MySQL] 2014/09/24 09:29:15 packets.go:92: write tcp 192.168.42.106:3306: broken pipe
[MySQL] 2014/09/24 09:29:15 statement.go:24: Invalid Connection
[MySQL] 2014/09/24 09:29:15 statement.go:24: Invalid Connection
[09/24/14 09:29:15] [EROR] db query: driver: bad connection
wakes: [March{uid:78, mid:110, type:attack, state:marching, (41, 47), due:2014-09-28 08:17:11 +0000 UTC} March{uid:78, mid:110, type:attack, state:marching, (41, 47), due:2014-09-28 08:17:11 +0000 UTC} March{uid:78, mid:110, type:attack, state:marching, (41, 47), due:2014-09-28 08:17:11 +0000 UTC} March{uid:78, mid:110, type:attack, state:marching, (41, 47), due:2014-09-28 08:17:11 +0000 UTC}]

PDDL

### TODO
*   metrics of php request handling
*   teleport fails need the target tile most recent info
*   Write/Read timeout and check N in loop
*   A@t0 build a 5s farm which arrives at php at t3, at t2 A's research timeout, what will happen if opTime/serverTime?
*   can a player send N marches to the same tile?
*   simulate mysql shutdown
    - done! golang mysql driver with breaker will handle this
*   WHERE UNIX_TIMESTAMP(time_end) index hit
    - need to optimize DB index
*   worker throttle
    - we can't have toooo many callbacks concurrently, use channel for throttle, easy...
*   FlightKey
    - Job flight is Job row or (uid, jobId)?
    - what if call back a Job but next wakeup turn the job status changed?
*   merge Job and March for a given player due time same into 1 callback
    - e,g. a HealTroop Job and attachee march arrives at the same time, which first run?
    - do we really need this?
    - maybe we dont care about this if we lock player correctly
*   handles NULL column
    - march.type done if it's NULL, what about others?
    - maybe we should let DB handle this
    - but mysql enum datatype can't handle this automatically
*   MaxRetries has bug
    - for a give tile, it will give up callback even after N success callback
