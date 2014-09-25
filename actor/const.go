package actor

const (
	JOB_QUERY = "SELECT uid,job_id,city_id,event_type,time_start,time_end,trace FROM Job WHERE unix_timestamp(time_end)<=? ORDER BY time_end ASC" // FOR UPDATE NOWAIT?

	MARCH_QUERY = "SELECT uid,march_id,end_x,end_y,state,end_time FROM March WHERE state!='done' AND unix_timestamp(end_time)<=? ORDER BY end_time ASC, end_x, end_y"

	PVE_QUERY = "SELECT uid,march_id,state,end_time FROM PveMarch WHERE unix_timestamp(end_time)<=? AND state!='done' ORDER BY end_time ASC"
)
