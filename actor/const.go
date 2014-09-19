package actor

const (
	CONTENT_TYPE_JSON = "application/json"

	RESPONSE_OK    = "ok"
	RESPONSE_RETRY = "retry"

	JOB_QUERY = "SELECT uid,job_id,time_end FROM Job WHERE time_end>=?"
)
