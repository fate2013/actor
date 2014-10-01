package actor

import (
	"database/sql"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type MysqlPoller struct {
	interval           time.Duration
	slowQueryThreshold time.Duration

	stopChan chan bool

	mysql   *sql.DB
	breaker *breaker.Consecutive

	jobQueryStmt   *sql.Stmt
	marchQueryStmt *sql.Stmt
	pveQueryStmt   *sql.Stmt

	latency metrics.Histogram
}

func NewMysqlPoller(interval time.Duration, slowQueryThreshold time.Duration,
	my *ConfigMysqlInstance, query *ConfigMysqlQuery, bc *ConfigBreaker) *MysqlPoller {
	this := new(MysqlPoller)
	this.interval = interval
	this.slowQueryThreshold = slowQueryThreshold

	this.stopChan = make(chan bool)

	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.db", this.latency)

	this.breaker = &breaker.Consecutive{
		FailureAllowance: bc.FailureAllowance,
		RetryTimeout:     bc.RetryTimeout}
	var err error
	this.mysql, err = sql.Open("mysql", my.DSN())
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

	this.jobQueryStmt, err = this.mysql.Prepare(query.Job)
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}

	this.marchQueryStmt, err = this.mysql.Prepare(query.March)
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}

	this.pveQueryStmt, err = this.mysql.Prepare(query.Pve)
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}

	return this
}

func (this *MysqlPoller) Poll(ch chan<- Wakeable) {
	go this.poll("job", ch)
	go this.poll("march", ch)
	go this.poll("pve", ch)

	<-this.stopChan

	this.jobQueryStmt.Close()
	this.marchQueryStmt.Close()
	this.pveQueryStmt.Close()
}

func (this *MysqlPoller) poll(typ string, ch chan<- Wakeable) {
	ticker := time.NewTicker(this.interval)
	defer ticker.Stop()

	var ws []Wakeable
	for now := range ticker.C {
		ws = this.fetchWakeables(typ, now)
		if len(ws) == 0 {
			continue
		}

		log.Debug("wakes[%s]: %+v", typ, ws)

		for _, w := range ws {
			ch <- w
		}
	}
}

func (this *MysqlPoller) fetchWakeables(typ string, dueTime time.Time) (ws []Wakeable) {
	ws = make([]Wakeable, 0, 100)
	var stmt *sql.Stmt
	switch typ {
	case "job":
		stmt = this.jobQueryStmt

	case "march":
		stmt = this.marchQueryStmt

	case "pve":
		stmt = this.pveQueryStmt
	}

	rows, err := this.Query(stmt, dueTime.Unix())
	if err != nil {
		log.Error("db query: %s", err.Error())

		return
	}

	this.latency.Update(time.Since(dueTime).Nanoseconds() / 1e6)

	switch typ {
	case "job":
		for rows.Next() {
			var w Job
			err = rows.Scan(&w.Uid, &w.JobId, &w.CityId,
				&w.Type, &w.TimeStart, &w.TimeEnd, &w.Trace)
			if err != nil {
				log.Error("db scan: %s", err.Error())
				continue
			}

			ws = append(ws, &w)
		}

	case "march":
		for rows.Next() {
			var w March
			err = rows.Scan(&w.Uid, &w.MarchId, &w.X1, &w.Y1, &w.Type,
				&w.State, &w.EndTime)
			if err != nil {
				log.Error("db scan: %s", err.Error())
				continue
			}

			ws = append(ws, &w)
		}

	case "pve":
		for rows.Next() {
			var w Pve
			err = rows.Scan(&w.Uid, &w.MarchId, &w.State, &w.EndTime)
			if err != nil {
				log.Error("db scan: %s", err.Error())
				continue
			}

			ws = append(ws, &w)
		}
	}

	rows.Close()
	return
}

func (this *MysqlPoller) Query(stmt *sql.Stmt,
	args ...interface{}) (rows *sql.Rows, err error) {
	//log.Debug("%+v, args=%+v", *stmt, args)

	if this.breaker.Open() {
		return nil, ErrCircuitOpen
	}

	t0 := time.Now()
	rows, err = stmt.Query(args...)
	if err != nil {
		this.breaker.Fail()
		return
	} else {
		this.breaker.Succeed()
	}

	elapsed := time.Since(t0)
	if elapsed > this.slowQueryThreshold {
		log.Warn("slow query:%v, %+v", elapsed, *stmt)
	}

	return
}

func (this *MysqlPoller) Stop() {
	close(this.stopChan)
}
