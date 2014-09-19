/*

                           +-- mysql
                           |
          persistent conn  |-- mysql
   actord -----------------|
     |           job poll  |-- mysql
     | http                |
     | callback            +-- mysql
     |
    php ---+
     |     | lock
     +-----+

actor schedules march event: arrive(speedup), recall(homeBack)


           client                client
             |                     |
          +-----+               +-----+
          |     |               |     |
       IM ^     V http       IM |     | http
          |     |               |     |
    500ms |     | 800ms         |     |
          +-----+               +-----+
             |                     |
             V consume             ^ pilliage
             |                     |        
             |---------------------+           
        sync |                                  sched
             |                                 +-----+
             |  10us          10us       5ms   |     | 1ms
            php ----- syslogng ---- proxy --- actor -+
             |  async                         |
             |                                |
             | resource                       |
             +-------------<------------------+
             50-800ms


                march
      +-----------------------------+
      |                             |
      |         lock inc, callback  |
    actor sched -----------------> php <---- unity3d
                                    |
                                    | check lock
                                    |
                                  city



RTS: synchronization is king

Synchronous execution

lock step

surrogate

serial scheduler
soft real time system
local latency and remote latency


predictable vs unpredictable events
predictable events are generated independently at each node, since each
node has the full game state

recoverable abort


A@B


request db token(permission) for each db operation

* delay req
* refuse req

concurrent issues: atomic, consistency, isolation, durability

hero points:
    deduct hero points
    (another req might arrive -- being attacked)  external actions
    add hero skill, which have effect in combat

build a farm:
    deduct resource
    (another req might arrive -- defeated and pillaged by attacker)

consume:
    deduct consumable
    (another req might arrive)
    add player resource

bank transfer:
    deduct A's account X
    (another req might arrive)
    add to B's account X


combat as a single request, it must be atomic and isolated, but its duration is 0ms-30s
    begin combat
    get troops and hero in both sides   --
    calculate combat
    get loser pillagings                --
    add pillagings to winner march
    commit combat

影响的是：
    战斗结果(输、赢和程度， 损失多少兵)
    战斗后战利品掠夺(掠夺时间内，不能有战利品相关数据变化)

critical sections(collision):
    1. tile encamp                                          - between marches
        try {
            tile.create
        } catch PDOException {
            // some other march already occpied this tile
            fire combat
        }

        or by actord, but what if 1s between A@C and B@C?
        ring buffer? actored need track slide window for each critical sections
        a, b attack m at t, c attack m at t+1?


    2. parties taking part in combat(a,b,c attack d)        - between marches
        who combat first, seond?

    3. pillage of attackee                                  - between march and attackee behavior
        3.1 resource
        3.2 hero capture

        crete global shard table combat {
            attacker_uid
            attackee_uid
            target_id (cityId | tileId) 
            ctime
            uniq key(attacker_uid, attackee_uid, target_id)
            key(attackee_uid)
        }

        try {
            c = combat.create // this must be the first write of db in this request
        } catch PDOException {
            // another combat for this attackee is happening
            return false
        } 
        doCombat
        pillaging and loot drop
        c.delete // what if this fail?

    4. a gathering, b,c arrives at roughly the same time(actor need serialize b, c combat?)
       b arrives at t, c arrives at t+1

       try {
           g = gathering.create
       } catch PDOException {
           // another player is within gathering request, but that request is not finished
           return false
       }

       at t, actor callback php to let (a, b) combat
       at t+1, active callback php to let (a, c) combat
       but the 2 callback might interleave, need better way to serialize this

    5. user consumes items/build construction/research/etc
       1. user.gold -= x
       2. city.resource -= y
       3. cityTile.save
       If combat begins right before 2, will lead to NotEnoughResourceException while user lost gold

       how to lock?
       if combat.exists(uid) {
           return false
       }

       combat.lock(uid)



lockings
========

lockUser(uid)  resource
lockTile(k, x, y)  encamp/gather


principle:
    php can't sleep and wait
    critical section can be held for long


sematics:
    enter crital section
    leave crital section



attacker need lock during combat?
    1. before attack, php will wake up attacker and attackee pending jobs, including boost
    2. attacker use hero points to boost hero | attacker consume item for boost related


attackee lock:
    1. 


problems:
1. callback to php combat while user consumes an item req arrives, 2 php instances race condition
2. 

root reason:
requests are async and sequence is unexpected


    lock -> combat -> unlock

    what about using mysql as mutex?
    encamp arrive: tile create conflict and fail(march1 arrive at t, march2 arrive at same tile at t+1,
    but march2 tile::create may be called first)


*/
package actor


