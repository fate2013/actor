package actor

import (
	"database/sql"
	log "github.com/funkygao/log4go"
	"time"
)

type poller struct {
	interval time.Duration
	mysql    *mysql
	jobStmt  *sql.Stmt
}

func newPoller(interval time.Duration, mysql *mysql) *poller {
	this := new(poller)
	this.interval = interval
	this.mysql = mysql
	var err error
	this.jobStmt, err = this.mysql.db.Prepare("SELECT uid,job_id,time_end FROM Job WHERE time_end>=?")
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}

	return this
}

func (this *poller) run(jobCh chan<- job) {
	ticker := time.NewTicker(this.interval)
	defer ticker.Stop()

	var job job
	for now := range ticker.C {
		rows, err := this.jobStmt.Query(now.Unix())
		if err != nil {
			log.Error("db query error: %s", err.Error())
			continue
		}

		for rows.Next() {
			err = rows.Scan(&job.Uid, &job.JobId, &job.dueTime)
			if err != nil {
				log.Error("db scan error: %s", err.Error())
				continue
			}

			log.Debug("due %+v", job)
			jobCh <- job
		}

	}

}

func (this *poller) stop() {
	this.jobStmt.Close()
}
