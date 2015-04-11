package httpq

type Queue interface {
	Push([]byte) error
	Pop() ([]byte, error)
	Size() (int64, error)
}
