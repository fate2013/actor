package actor

const (
	RESPONSE_OK    = 1
	RESPONSE_RETRY = 5

	JOB_QUERY   = "SELECT uid,job_id,city_id,event_type,time_start,time_end,trace FROM Job WHERE unix_timestamp(time_end)<=? ORDER BY time_end ASC" // FOR UPDATE NOWAIT?
	JOB_KILL    = "DELETE FROM Job WHERE uid=? AND job_id=?"                                                                                        // TODO
	MARCH_QUERY = "SELECT uid,march_id,end_x,end_y,state,end_time FROM March WHERE state!='done' AND unix_timestamp(end_time)<=? ORDER BY end_time ASC, end_x, end_y"
)
