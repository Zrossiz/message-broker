package handler

import (
	"encoding/json"
	"io"
	"net"
	"time"

	"github.com/Zrossiz/message-broker/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Request struct {
	Action  string `json:"action"`
	Message []byte `json:"message"`
	Offset  uint64 `json:"offset"`
}

type Handler struct {
	srv MessageService
}

func New(srv MessageService) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) HandleMessage(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var req Request

		err := decoder.Decode(&req)
		if err != nil {
			if err != io.EOF {
				zap.L().Error("decode request", zap.Error(err))
			}
		}

		switch req.Action {
		case "produce":
			msg := model.Message{
				ID:        uuid.New(),
				Message:   req.Message,
				CreatedAt: time.Now(),
			}

			err = h.srv.Append(msg)
			if err != nil {
				zap.L().Error("append err", zap.Error(err))
				return
			}

		case "fetch":
			msg, err := h.srv.Get(req.Offset)
			if err != nil {
				zap.L().Error("get message err", zap.Error(err))
				return
			}

			err = encoder.Encode(msg)
			if err != nil {
				zap.L().Error("encode message err", zap.Error(err))
				return
			}
		default:
			zap.L().Error("invalid action", zap.Error(err))
			return
		}
	}
}
