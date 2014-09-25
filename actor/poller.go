package actor

type Poller interface {
	Run(ch chan<- Wakeable)
	Stop()
}
