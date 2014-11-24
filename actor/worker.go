package actor

type Worker interface {
	Start()
	Wake(w Wakeable)
	FlightCount() int
	Flights() map[string]interface{}
}
