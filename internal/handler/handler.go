package handler

import (
	"context"
	"encoding/json"
	"io"
	"net"

	"go.uber.org/zap"
)

type ResMessage struct {
	Success bool `json:"success"`
}

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

func (h *Handler) HandleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var req Request

			err := decoder.Decode(&req)
			if err != nil {
				if err != io.EOF {
					zap.L().Error("decode request", zap.Error(err))
				}
			}

			switch req.Action {
			case "produce":
				h.handleProduce(req, encoder)
			case "fetch":
				h.handleFetch(req, encoder)
			default:
				zap.L().Error("invalid action", zap.Error(err))
				return
			}
		}
	}
}

func (h *Handler) handleFetch(req Request, enc *json.Encoder) {
	msg, err := h.srv.Get(req.Offset)
	if err != nil {
		zap.L().Error("get message err", zap.Error(err))
		return
	}

	err = enc.Encode(msg)
	if err != nil {
		zap.L().Error("encode message err", zap.Error(err))
		return
	}
}

func (h *Handler) handleProduce(req Request, enc *json.Encoder) {
	err := h.srv.Append(req.Message)
	if err != nil {
		zap.L().Error("append err", zap.Error(err))
		return
	}

	res := ResMessage{
		Success: true,
	}

	err = enc.Encode(res)
	if err != nil {
		zap.L().Error("encode response err", zap.Error(err))
		return
	}
}
