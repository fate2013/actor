/*

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
            c = combat.create
        } catch PDOException {
            // another combat for this attackee is happening
            return false
        } 
        doCombat
        pillaging and loot drop
        c.delete // what if this fail?

    4. a gathering, b,c arrives at roughly the same time(actor need serialize b, c combat?)
       b arrives at t, c arrives at t+1

       at t, actor callback php to let (a, b) combat
       at t+1, active callback php to let (a, c) combat
       but the 2 callback might interleave, need better way to serialize this

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

