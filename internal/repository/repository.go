package repository

type Repository struct {
	Message MessageRepo
}

func New(filepath string) *Repository {
	return &Repository{
		Message: *NewMessage(filepath),
	}
}
