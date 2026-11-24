package handler

type MessageService interface {
	Append(msg []byte) error
	Get(offset uint64) ([]byte, error)
}
