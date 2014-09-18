package queue

var (
	keyFuncWhen GreaterThanFunc = func(l, r interface{}) bool {
		return l.(int64) > r.(int64)
	}
)

type GreaterThanFunc func(l, r interface{}) bool

func (f GreaterThanFunc) Descending() bool {
	return false
}

func (f GreaterThanFunc) Compare(l, r interface{}) bool {
	return f(l, r)
}
