package actor

import (
	"database/sql"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"time"
)

type MysqlPoller struct {
	interval time.Duration
	mysql    *mysql
	jobStmt  *sql.Stmt
	latency  metrics.Histogram
}

func NewMysqlPoller(interval time.Duration, my *ConfigMysqlInstance,
	breaker *ConfigBreaker) *MysqlPoller {
	this := new(MysqlPoller)
	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.db", this.latency)
	this.interval = interval

	this.mysql = newMysql(my.DSN(), breaker)
	err := this.mysql.Open()
	if err != nil {
		log.Critical("open mysql[%+v] failed: %s", *my, err)
		return nil
	}

	this.jobStmt, err = this.mysql.db.Prepare("SELECT uid,job_id,time_end FROM Job WHERE time_end>=?")
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}

	return this
}

func (this *MysqlPoller) Run(jobCh chan<- Job) {
	ticker := time.NewTicker(this.interval)
	defer ticker.Stop()

	var (
		job Job
		t0  time.Time
	)
	for now := range ticker.C {
		t0 = time.Now()
		rows, err := this.jobStmt.Query(now.Unix())
		if err != nil {
			log.Error("db query error: %s", err.Error())
			continue
		}
		this.latency.Update(time.Since(t0).Nanoseconds() / 1e6)

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

func (this *MysqlPoller) Stop() {
	this.jobStmt.Close()
}
