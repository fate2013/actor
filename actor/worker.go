package actor

type Worker interface {
	Wake(w Wakeable) (retry bool)
	FlightCount() int
	Flights() map[string]interface{}
}
