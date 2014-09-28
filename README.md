DragonServer
============

RTS game scheduler to handle the following challanges:
* concurrency isolation
* scheduling of delayed jobs

### Terms
*   hit
    - a scheduling interval: 1s

### TODO
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

### Benefits of introducing actor
*   actor make multiple user concurrency problem into 2-1 user problem
*   the only instant multiple user problem is: teleport(with help of DB integrity it gets solved)
    - it has problem because of (instant vs job conflicts)

### Assumptions
*   concurrency(time) accuracy is in second
    - there is no absolute order although in real world we can't step into the same rive at the same time
    - it means if 2-more events happens at the same second, they are concurrent
    - concurrent events will be serialized by actor or by chance(lock, depends on who accquired the lock first)
*   combat is atomic(always A@B)
    - combat must lock attacker & attackee

### Problems
*   server push lost and out-of-order
*   race condition between instant action and job based action

### Lock, lock what?
*   just user lock: current user(instant action)
*   attacker and attachee
*   any action may fail because of the mutex lock
    - job/march fail: actor will retry after 1s
    - instant action fail: client will popup dialog to inform user: user decide when to retry

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
