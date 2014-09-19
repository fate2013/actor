package actor

type Callbacker interface {
	Call(j Job)
}
