package actor

const (
	RESPONSE_OK    = 1
	RESPONSE_RETRY = 5

	JOB_QUERY = "SELECT uid,job_id,city_id,event_type,time_start,time_end,trace FROM Job WHERE time_end<=NOW() ORDER BY time_end ASC" // FOR UPDATE NOWAIT?
	JOB_KILL  = "DELETE FROM Job WHERE uid=? AND job_id=?"                                                                            // TODO
)
