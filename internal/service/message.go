package service

import (
	"github.com/Zrossiz/message-broker/internal/model"
)

type MessageService struct {
	messageRepo MessageRepo
}

func NewMessageService(messageRepo MessageRepo) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

func (s *MessageService) Append(msg model.Message) error {
	return s.messageRepo.Append(msg)
}

func (s *MessageService) Get(offset uint64) (*model.Message, error) {
	return s.messageRepo.Get(offset)
}
