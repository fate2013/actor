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
*/
package actor
