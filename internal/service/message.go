package service

type MessageService struct {
	messageRepo MessageRepo
}

func NewMessageService(messageRepo MessageRepo) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

func (s *MessageService) Append(msg []byte) error {
	return s.messageRepo.Append(msg)
}

func (s *MessageService) Get(offset uint64) ([]byte, error) {
	return s.messageRepo.Get(offset)
}
