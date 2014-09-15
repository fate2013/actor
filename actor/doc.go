/*
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

*/
package actor


