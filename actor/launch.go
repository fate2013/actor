package actor

func (this *Actor) ServeForever() {
	this.launchHttpServ()
	defer this.stopHttpServ()

	go this.runAcceptor()

	this.runScheduler()
}
