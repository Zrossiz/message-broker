package service

type Service struct {
	Message MessageRepo
}

func New(messageRepo MessageRepo) *Service {
	return &Service{
		Message: NewMessageService(messageRepo),
	}
}
