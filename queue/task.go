package queue

type Task struct {
	Who   int64
	When  int64
	Where int
	What  string
}
