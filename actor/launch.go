package actor

func (this *Actor) ServeForever() {
	this.launchHttpServ()
	defer this.stopHttpServ()

	this.replicator.Replay()
	this.replicator.Start()

	go this.runAcceptor()

	this.runScheduler()
}
