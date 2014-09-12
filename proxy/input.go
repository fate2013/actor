package proxy

type Input interface {
	Reader() chan []byte
}
