package actor

type Poller interface {
	Run(jCh chan<- Job, mCh chan<- March, pCh chan<- Pve)
	Stop()
}
