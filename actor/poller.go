package actor

type Poller interface {
	Run(ch chan<- Schedulable)
	Stop()
}
