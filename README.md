DragonServer
============

### TODO
*   simulate mysql shutdown
*   WHERE UNIX_TIMESTAMP(time_end) index hit
*   worker throttle
*   FlightKey


[MySQL] 2014/09/24 09:29:15 packets.go:30: read tcp 192.168.42.106:3306: operation timed out
[MySQL] 2014/09/24 09:29:15 packets.go:92: write tcp 192.168.42.106:3306: broken pipe
[MySQL] 2014/09/24 09:29:15 statement.go:24: Invalid Connection
[MySQL] 2014/09/24 09:29:15 statement.go:24: Invalid Connection
[09/24/14 09:29:15] [EROR] db query: driver: bad connection



teleport race with actor


pve march callback, then the march state changed, lock fails, -> 2 concurrent callbacks for the 1 march

Consider, for example, production problems in which a final product is built from raw materials 
called resources. 
In such domains there can be dependencies between actions: some actions may accumulate 
certain resources, while other actions consume resources to produce something. 


RTS objective being to achieve military or territorial su- periority over other players or the computer.
Central to RTS game-play are two key problem domains, resource produc- tion and tactical battles.

server-client instead of client-server arch

                system
                  |
            +--  tile    --+
            |              |
            |              |
            +-- units    --+
   player --|-- resource --|-- opponent
            +-- hero xp  --+


actions:
    consume
    research
    boost
    construct
    


    (precondition, effect) - atomic pair

player: client, server, actor, batch scripts,

                   swoole    +-- player
              +--- php-fpm --|
              |              +-- actor
              |
        db ---|--- mq
              |
              |
              +--- batch script(resource tile generator)


The common objective in RTS games is to eliminate other players through military superiority.
Players first instruct workers to gather resources, then use those resources to build more workers 
and structures that can create military units. These are then sent to battle the enemy in real-time.


this record leads to deadloop, php always returns marching, never change it

                   uid: 43
              march_id: 70
               city_id: 43
               opp_uid: 16
              world_id: 1
                  type: NULL
               start_x: 20
               start_y: 28
                 end_x: 22
                 end_y: 22
            start_time: 2014-09-26 23:09:38
              end_time: 2014-09-26 23:09:44
                 speed: 1
          start_troops: {"cavalry_t1":"1"}
                troops: {"cavalry_t1":"1"}
      opp_start_troops: []
            opp_troops: []
troops_sentto_hospital: []
             resources: []
                 ctime: 2014-09-26 23:09:38
                 mtime: 2014-09-26 23:09:39
                 state: marching
            battle_log: []
