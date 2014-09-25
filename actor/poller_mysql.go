package actor

import (
	"database/sql"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"time"
)

type MysqlPoller struct {
	interval          time.Duration
	mysql             *mysql
	jobQueryStmt      *sql.Stmt
	marchQueryStmt    *sql.Stmt
	jobQueryLatency   metrics.Histogram
	marchQueryLatency metrics.Histogram
}

func NewMysqlPoller(interval time.Duration, my *ConfigMysqlInstance,
	breaker *ConfigBreaker) *MysqlPoller {
	this := new(MysqlPoller)
	this.interval = interval

	this.jobQueryLatency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.db.job", this.jobQueryLatency)
	this.marchQueryLatency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.db.march", this.marchQueryLatency)

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

	this.jobQueryStmt, err = this.mysql.db.Prepare(JOB_QUERY)
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}

	this.marchQueryStmt, err = this.mysql.db.Prepare(MARCH_QUERY)
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}

	return this
}

// TODO select timeout jobs, then delete them
// in case of multiple actord, check delete afftectedRows==rowsCount, then dispatch job
func (this *MysqlPoller) Run(jobCh chan<- Job, marchChan chan<- March) {
	defer this.Stop()

	this.pollMarch(marchChan)

	//this.pollJob(jobCh)
}

func (this *MysqlPoller) pollJob(jobCh chan<- Job) {
	ticker := time.NewTicker(this.interval)
	defer ticker.Stop()

	var jobs []Job
	for now := range ticker.C {
		jobs = this.fetchReadyJobs(now)
		if len(jobs) == 0 {
			continue
		}

		log.Debug("waking up %+v", jobs)

		for _, job := range jobs {
			jobCh <- job
		}

	}
}

func (this *MysqlPoller) pollMarch(marchCh chan<- March) {
	ticker := time.NewTicker(this.interval)
	defer ticker.Stop()

	var marches MarchGroup
	for now := range ticker.C {
		marches = this.fetchReadyMarches(now)
		if len(marches) > 0 {
			log.Debug("due %+v", marches)
		}

		for _, march := range marches {
			marchCh <- march
		}
	}

}

func (this *MysqlPoller) fetchReadyMarches(dueTime time.Time) (marches MarchGroup) {
	rows, err := this.marchQueryStmt.Query(dueTime.Unix())
	if err != nil {
		log.Error("db query: %s", err.Error())
		return
	}

	this.marchQueryLatency.Update(time.Since(dueTime).Nanoseconds() / 1e6)

	var march March
	for rows.Next() {
		err = rows.Scan(&march.Uid, &march.MarchId, &march.X1, &march.Y1,
			&march.State, &march.EndTime)
		if err != nil {
			log.Error("db scan: %s", err.Error())
			continue
		}

		marches = append(marches, march)
	}

	//marches.sortByDestination()

	rows.Close()
	return
}

// TODO disable autocommit
// Job rows: time_end with 1, 5, 9
// select * from Job where time_end<=8
// php(player) speedup 9 to 4
// delete from Job where time_end<=8 will miss job(4)
// contention exists between actord and php(because job can pause/speedup/cancel)
func (this *MysqlPoller) fetchReadyJobs(dueTime time.Time) (jobs []Job) {
	//jobs = make([]Job, 0, 100)

	rows, err := this.jobQueryStmt.Query(dueTime.Unix())
	if err != nil {
		log.Error("db query: %s", err.Error())
		return
	}

	this.jobQueryLatency.Update(time.Since(dueTime).Nanoseconds() / 1e6)

	var job Job
	for rows.Next() {
		err = rows.Scan(&job.Uid, &job.JobId, &job.CityId,
			&job.Type, &job.TimeStart, &job.TimeEnd, &job.Trace)
		if err != nil {
			log.Error("db scan: %s", err.Error())
			continue
		}

		jobs = append(jobs, job)
	}

	rows.Close() // Failing to use rows.Close() or stmt.Close() can cause exhaustion of resources

	return
}

func (this *MysqlPoller) Stop() {
	this.jobQueryStmt.Close()
	this.marchQueryStmt.Close()
}
