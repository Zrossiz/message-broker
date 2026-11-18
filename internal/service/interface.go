package service

import (
	"github.com/Zrossiz/message-broker/internal/model"
)

type MessageRepo interface {
	Append(msg model.Message) error
	Get(offset uint64) (*model.Message, error)
}
