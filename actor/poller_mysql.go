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
		log.Critical("ping mysql[%s]: %s", my.DSN(), err)
		return nil
	}
	log.Debug("mysql connected: %s", my.DSN())

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
		rows, err := this.queryStmt.Query()
		if err != nil {
			log.Error("db query error: %s", err.Error())
			continue
		}
		this.queryLatency.Update(time.Since(now).Nanoseconds() / 1e6)

		for rows.Next() {
			err = rows.Scan(&job.Uid, &job.JobId, &job.CityId,
				&job.Type, &job.TimeStart, &job.TimeEnd, &job.Trace)
			if err != nil {
				log.Error("db scan error: %s", err.Error())
				continue
			}

			log.Debug("waking up %s", job)
			jobCh <- job
		}

		rows.Close() // Failing to use rows.Close() or stmt.Close() can cause exhaustion of resources
	}

}

// TODO disable autocommit
// Job rows: time_end with 1, 5, 9
// select * from Job where time_end<=8
// php(player) speedup 9 to 4
// delete from Job where time_end<=8 will miss job(4)
// contention exists between actord and php(because job can pause/speedup/cancel)
func (this *MysqlPoller) fetchReadyJobs(dueTime time.Time) (jobs []Job) {
	jobs = make([]Job, 0, 100)
	var job Job

	rows, err := this.queryStmt.Query()
	if err != nil {
		log.Error("db query error: %s", err.Error())
		return
	}

	this.queryLatency.Update(time.Since(dueTime).Nanoseconds() / 1e6)

	for rows.Next() {
		err = rows.Scan(&job.Uid, &job.JobId, &job.CityId,
			&job.Type, &job.TimeStart, &job.TimeEnd, &job.Trace)
		if err != nil {
			log.Error("db scan error: %s", err.Error())
			continue
		}

		res, err := this.killStmt.Exec(job.Uid, job.JobId)
		if err != nil {
			log.Error("kill job[%+v]: %s", job, err)
			continue
		}
		if n, _ := res.RowsAffected(); n != 1 {
			// another process has killed this job since I query
			log.Warn("job killed by another instance: %+v", job)
			continue
		}

		log.Debug("job[%+v] killed", job)
		jobs = append(jobs, job)
	}

	log.Debug("due jobs: %+v", jobs)

	return
}

func (this *MysqlPoller) Stop() {
	this.queryStmt.Close()
	this.killStmt.Close()
}
