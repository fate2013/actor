package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

type poller struct {
	actor *Actor
	mysql *mysql
}

func newPoller(actor *Actor, mysql *mysql) *poller {
	this := new(poller)
	this.actor = actor
	this.mysql = mysql
	return this

}

func (this *poller) run(jobCh chan<- string) {
	ticker := time.NewTicker(this.actor.config.ScheduleInterval)
	defer ticker.Stop()

	for now := range ticker.C {
		rows, err := this.mysql.Query("SELECT uid FROM Job WHERE time_end>=", now.Unix())
		if err != nil {
			log.Error("db error: %s", err.Error())
			continue
		}

		cols, _ := rows.Columns()
		colsN := len(cols)
		values := make([]interface{}, colsN)
		valuePtrs := make([]interface{}, colsN)
		for rows.Next() {
			for i, _ := range cols {
				valuePtrs[i] = &values[i]
				rows.Scan(valuePtrs...)
				uid := values[0].(int64)
			}
		}
	}

}

func (this *poller) lock() {
	// use memcache.add for distributed atomic lock
}
