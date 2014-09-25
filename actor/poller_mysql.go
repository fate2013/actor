package actor

import (
	"database/sql"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"time"
)

type MysqlPoller struct {
	interval time.Duration
	stopChan chan bool

	breaker *breaker.Consecutive

	mysql          *mysql
	jobQueryStmt   *sql.Stmt
	marchQueryStmt *sql.Stmt
	pveQueryStmt   *sql.Stmt

	latency metrics.Histogram
}

func NewMysqlPoller(interval time.Duration, my *ConfigMysqlInstance,
	bc *ConfigBreaker) *MysqlPoller {
	this := new(MysqlPoller)
	this.interval = interval
	this.stopChan = make(chan bool)

	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.db", this.latency)

	this.mysql = newMysql(my.DSN(), bc)
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

	this.pveQueryStmt, err = this.mysql.db.Prepare(PVE_QUERY)
	if err != nil {
		log.Critical("db prepare err: %s", err.Error())
		return nil
	}

	this.breaker = &breaker.Consecutive{
		FailureAllowance: bc.FailureAllowance,
		RetryTimeout:     bc.RetryTimeout}

	return this
}

func (this *MysqlPoller) Run(ch chan<- Schedulable) {
	go this.poll("job", ch)
	go this.poll("march", ch)
	go this.poll("pve", ch)

	<-this.stopChan

	this.jobQueryStmt.Close()
	this.marchQueryStmt.Close()
	this.pveQueryStmt.Close()
}

func (this *MysqlPoller) poll(typ string, ch chan<- Schedulable) {
	ticker := time.NewTicker(this.interval)
	defer ticker.Stop()

	var ss []Schedulable
	for now := range ticker.C {
		ss = this.fetchSchedulables(typ, now)
		if len(ss) == 0 {
			continue
		}

		log.Debug("waking up %+v", ss)

		for _, s := range ss {
			ch <- s
		}
	}
}

func (this *MysqlPoller) fetchSchedulables(typ string, dueTime time.Time) (ss []Schedulable) {
	if this.breaker.Open() {
		log.Warn("breaker open %+v", *this.breaker)
		return
	}

	var stmt *sql.Stmt
	switch typ {
	case "job":
		stmt = this.jobQueryStmt

	case "march":
		stmt = this.marchQueryStmt

	case "pve":
		stmt = this.pveQueryStmt
	}

	rows, err := stmt.Query(dueTime.Unix())
	if err != nil {
		log.Error("db query: %s", err.Error())

		this.breaker.Fail()
		return
	} else {
		this.breaker.Succeed()
	}

	this.latency.Update(time.Since(dueTime).Nanoseconds() / 1e6)

	switch typ {
	case "job":
		var s Job
		for rows.Next() {
			err = rows.Scan(&s.Uid, &s.JobId, &s.CityId,
				&s.Type, &s.TimeStart, &s.TimeEnd, &s.Trace)
			if err != nil {
				log.Error("db scan: %s", err.Error())
				continue
			}

			ss = append(ss, &s)
		}

	case "march":
		var s March
		for rows.Next() {
			err = rows.Scan(&s.Uid, &s.MarchId, &s.X1, &s.Y1,
				&s.State, &s.EndTime)
			if err != nil {
				log.Error("db scan: %s", err.Error())
				continue
			}

			ss = append(ss, &s)
		}

	case "pve":
		var s Pve
		for rows.Next() {
			err = rows.Scan(&s.Uid, &s.MarchId, &s.State, &s.EndTime)
			if err != nil {
				log.Error("db scan: %s", err.Error())
				continue
			}

			ss = append(ss, &s)
		}
	}

	rows.Close()
	return
}

func (this *MysqlPoller) Stop() {
	close(this.stopChan)
}
