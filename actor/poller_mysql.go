package actor

import (
	"database/sql"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"time"
)

type MysqlPoller struct {
	interval     time.Duration
	mysql        *mysql
	queryStmt    *sql.Stmt
	killStmt     *sql.Stmt
	queryLatency metrics.Histogram
	killLatency  metrics.Histogram
}

func NewMysqlPoller(interval time.Duration, my *ConfigMysqlInstance,
	breaker *ConfigBreaker) *MysqlPoller {
	this := new(MysqlPoller)
	this.interval = interval
	this.queryLatency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.db.query", this.queryLatency)
	this.killLatency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.db.kill", this.killLatency)

	this.mysql = newMysql(my.DSN(), breaker)
	err := this.mysql.Open()
	if err != nil {
		log.Critical("open mysql[%+v] failed: %s", *my, err)
		return nil
	}
	err = this.mysql.Ping()
	if err != nil {
		log.Critical("ping mysql[%+v]: %s", *my, err)
		return nil
	}

	this.queryStmt, err = this.mysql.db.Prepare(JOB_QUERY)
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}
	this.killStmt, err = this.mysql.db.Prepare(JOB_KILL)
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}

	return this
}

// TODO select timeout jobs, then delete them
// in case of multiple actord, check delete afftectedRows==rowsCount, then dispatch job
func (this *MysqlPoller) Run(jobCh chan<- Job) {
	ticker := time.NewTicker(this.interval)
	defer func() {
		ticker.Stop()
		this.Stop()
	}()

	var job Job
	for now := range ticker.C {
		rows, err := this.queryStmt.Query(now.Unix())
		if err != nil {
			log.Error("db query error: %s", err.Error())
			continue
		}
		this.queryLatency.Update(time.Since(now).Nanoseconds() / 1e6)

		for rows.Next() {
			err = rows.Scan(&job.Uid, &job.JobId, &job.dueTime)
			if err != nil {
				log.Error("db scan error: %s", err.Error())
				continue
			}

			log.Debug("waking up %+v", job)
			jobCh <- job
		}
	}

}

func (this *MysqlPoller) fetchReadyJobs(dueTime time.Time) {
	//t0 := time.Now()

}

func (this *MysqlPoller) Stop() {
	this.queryStmt.Close()
	this.killStmt.Close()
}
