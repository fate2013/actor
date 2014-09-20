package actor

const (
	RESPONSE_OK    = 1
	RESPONSE_RETRY = 5

	JOB_QUERY = "SELECT uid,job_id,time_end FROM Job WHERE time_end<=? ORDER BY time_end ASC" // FOR UPDATE NOWAIT?
	JOB_KILL  = "DELETE FROM Job where time_end>=?"                                           // TODO
)
