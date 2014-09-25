package actor

type Callbacker interface {
	Call(s Schedulable) (retry bool)
}
