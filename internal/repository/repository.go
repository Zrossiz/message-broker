package repository

type Repository struct {
	Message *MessageRepo
}

func New(connectionPool *ConnectionPool) (*Repository, error) {
	return &Repository{
		Message: NewMessage(connectionPool.File),
	}, nil
}
