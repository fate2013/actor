package actor

type Worker interface {
	Wake(w Wakeable)
	FlightCount() int
	Flights() map[string]interface{}
}
