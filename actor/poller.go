package actor

type Poller interface {
	Run(ch chan<- Job)
	Stop()
}
