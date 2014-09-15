/*
actor schedules march event: arrive(speedup), recall(homeBack)

critical sections:
. encamp same tile at roughly the same time
. attack same city at roughly the same time
. arrive at gathering tile at roughly the same time
. wonder? mini wonder?
. TODO see march state diagram

all the critical sections share the same attribute:
same destination(geohash) at roughly the same time

problems:
1. callback to php combat while user consumes an item req arrives, 2 php instances race condition
2. 

root reason:
requests are async and sequence is unexpected

solutions:
how to turn async into sync

lock for sync(deadlock)


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

    lock -> combat -> unlock


*/
package actor


