package actor

type Worker interface {
	Wake(w Wakeable) (retry bool)
	InFlight() int
}
