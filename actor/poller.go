package actor

type poller struct {
	mysql *mysql
}

func newPoller(mysql *mysql) *poller {
	this := new(poller)
	this.mysql = mysql
	return this

}

func (this *poller) start(jobCh chan<- string) {
	jobCh <- "x"

}
