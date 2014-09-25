package actor

type Worker interface {
	Wake(w Wakeable) (retry bool)
	Outstandings() int
}
