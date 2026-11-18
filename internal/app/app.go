package app

import (
	"fmt"
	"net"

	"github.com/Zrossiz/message-broker/internal/config"
	"github.com/Zrossiz/message-broker/internal/handler"
	"github.com/Zrossiz/message-broker/internal/repository"
	"github.com/Zrossiz/message-broker/internal/service"
	"github.com/Zrossiz/message-broker/pkg/logger"
	"go.uber.org/zap"
)

func New() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("init config err: %w", err)
	}

	log, err := logger.New(cfg.App.LoggerLevel)
	if err != nil {
		return fmt.Errorf("init logger err: %w", err)
	}
	zap.ReplaceGlobals(log)
	zap.L().Info("logger inited")

	repo := repository.New(cfg.Repo.Filepath)
	zap.L().Info("repository inited")
	srv := service.New(&repo.Message)
	zap.L().Info("service inited")
	h := handler.New(srv.Message)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.App.Port))
	if err != nil {
		return fmt.Errorf("")
	}
	defer listener.Close()
	zap.L().Info("listener inited")

	zap.L().Info("start accepting tcp connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			zap.L().Error("accept conn err", zap.Error(err))
			continue
		}

		go h.HandleMessage(conn)
	}

}
