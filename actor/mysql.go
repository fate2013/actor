package actor

import (
	"database/sql"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	_ "github.com/go-sql-driver/mysql"
)

// A mysql conn to a single mysql instance
// Conn pool is natively supported by golang
type mysql struct {
	dsn     string
	db      *sql.DB
	breaker *breaker.Consecutive
}

func newMysql(dsn string, bc *ConfigBreaker) *mysql {
	this := new(mysql)
	this.dsn = dsn
	this.breaker = &breaker.Consecutive{
		FailureAllowance: bc.FailureAllowance,
		RetryTimeout:     bc.RetryTimeout}

	return this
}

func (this *mysql) Open() (err error) {
	this.db, err = sql.Open("mysql", this.dsn)
	return
}

func (this *mysql) Ping() error {
	if this.db == nil {
		return ErrNotOpen
	}

	return this.db.Ping()
}

func (this mysql) String() string {
	return this.dsn
}

func (this *mysql) Query(query string, args ...interface{}) (rows *sql.Rows,
	err error) {
	log.Debug("%s, args=%+v", query, args)

	if this.breaker.Open() {
		return nil, ErrCircuitOpen
	}

	rows, err = this.db.Query(query, args...)
	if err != nil {
		this.breaker.Fail()
	} else {
		this.breaker.Succeed()
	}

	return
}
