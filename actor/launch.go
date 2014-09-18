package actor

func (this *Actor) ServeForever() {
	this.launchHttpServ()
	defer this.stopHttpServ()

	this.runScheduler()
}
