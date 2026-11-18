package handler

import (
	"github.com/Zrossiz/message-broker/internal/model"
)

type MessageService interface {
	Append(msg model.Message) error
	Get(offset uint64) (*model.Message, error)
}
