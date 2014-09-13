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
1. callback to php combat while user consumes an item, 2 php instances race condition
2.


           client                client
             |                     |
          +-----+               +-----+
          |     |               |     |
       IM ^     V http       IM |     | http
          |     |               |     |
    500ms |     | 800ms         |     |
          +-----+               +-----+
             |                     |
             |---------------------+           sched
             |                                +-----+
             |  10us          10us       5ms  |     | 1ms
            php ---- syslogng ---- proxy --- actor -+
             |                                |
             |                                |
             +-------------<------------------+
             50-800ms
*/
package actor
