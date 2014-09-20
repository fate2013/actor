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

	tx, err := this.mysql.db.Begin()
	if err != nil {
		log.Critical("tx: %s", err.Error())
		return
	}
	txQueryStmt := tx.Stmt(this.queryStmt)
	txKillStmt := tx.Stmt(this.killStmt)
	defer func() {
		tx.Commit()
		txQueryStmt.Close()
		txKillStmt.Close()
	}()

	rows, err := txQueryStmt.Query(dueTime.Unix())
	if err != nil {
		log.Critical("query: %s", err.Error())
		return
	}

	t1 := time.Now() // query done
	this.queryLatency.Update(t1.Sub(dueTime).Nanoseconds() / 1e6)

	// marshal db rows to Job
	for rows.Next() {
		rows.Scan(&job.Uid, &job.JobId, &job.dueTime)

		jobs = append(jobs, job)
	}

	// kill the job to avoid being waken up in next round
	result, err := txKillStmt.Exec(dueTime.Unix())
	this.killLatency.Update(time.Since(t1).Nanoseconds() / 1e6)

	if n, _ := result.RowsAffected(); n > 0 {
		log.Debug("%d jobs killed", n)
	}

	return
}

func (this *MysqlPoller) Stop() {
	this.queryStmt.Close()
	this.killStmt.Close()
}
