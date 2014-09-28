DragonServer
============

### TODO
*   simulate mysql shutdown
*   WHERE UNIX_TIMESTAMP(time_end) index hit
*   worker throttle
*   FlightKey
*   merge Job and March for a given player due time same into 1 callback
*   handles NULL column


[MySQL] 2014/09/24 09:29:15 packets.go:30: read tcp 192.168.42.106:3306: operation timed out
[MySQL] 2014/09/24 09:29:15 packets.go:92: write tcp 192.168.42.106:3306: broken pipe
[MySQL] 2014/09/24 09:29:15 statement.go:24: Invalid Connection
[MySQL] 2014/09/24 09:29:15 statement.go:24: Invalid Connection
[09/24/14 09:29:15] [EROR] db query: driver: bad connection
wakes: [March{uid:78, mid:110, type:attack, state:marching, (41, 47), due:2014-09-28 08:17:11 +0000 UTC} March{uid:78, mid:110, type:attack, state:marching, (41, 47), due:2014-09-28 08:17:11 +0000 UTC} March{uid:78, mid:110, type:attack, state:marching, (41, 47), due:2014-09-28 08:17:11 +0000 UTC} March{uid:78, mid:110, type:attack, state:marching, (41, 47), due:2014-09-28 08:17:11 +0000 UTC}]


(precondition, effect) - atomic pair

teleport race with actor

pve march callback, then the march state changed, lock fails, -> 2 concurrent callbacks for the 1 march

Job and March for a player wakes at the same time, should be merged into one callback
e,g HealTroop job and opponent march arrives at the same time?

    
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


