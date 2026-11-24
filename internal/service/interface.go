package service

type MessageRepo interface {
	Append([]byte) error
	Get(offset uint64) ([]byte, error)
}
