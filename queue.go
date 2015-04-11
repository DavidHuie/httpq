package httpq

type Queue interface {
	Push([]byte) error
	Pop() ([]byte, error)
	Size() (uint64, error)
}
