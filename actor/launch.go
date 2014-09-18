package actor

func (this *Actor) ServeForever() {
	this.stats.launchHttpServ()
	defer this.stats.stopHttpServ()

	this.runScheduler()
}
