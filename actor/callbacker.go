package actor

type Callbacker interface {
	Call(j Job) (retry bool)
	Play(m March) (retry bool)
	Pve(p Pve) (retry bool)
}
