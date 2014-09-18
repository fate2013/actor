package actor

func (this *Actor) ServeForever() {
	go this.scheduler.run()

	this.stats.run()
}
