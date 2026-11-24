package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zrossiz/message-broker/internal/config"
	"github.com/Zrossiz/message-broker/internal/handler"
	"github.com/Zrossiz/message-broker/internal/repository"
	"github.com/Zrossiz/message-broker/internal/service"
	"github.com/Zrossiz/message-broker/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	server net.Listener
	ctx    context.Context

	cfg      *config.Config
	handler  *handler.Handler
	connPool *repository.ConnectionPool
}

// TODO: сделать грейсфул шатдаун
// 1. Закрывать файл
// 2. Блокировать возможность подключения для новых соединений
// 3. Заканчивать обрабатывать текущие
func New() (*App, error) {
	var app App

	app.ctx, _ = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("init config err: %w", err)
	}
	app.cfg = cfg

	log, err := logger.New(cfg.App.LoggerLevel)
	if err != nil {
		return nil, fmt.Errorf("init logger err: %w", err)
	}
	zap.ReplaceGlobals(log)
	zap.L().Info("logger inited")

	app.connPool, err = repository.NewConnectionPool(cfg.Repo)
	if err != nil {
		return nil, err
	}

	repo, err := repository.New(app.connPool)
	if err != nil {
		return nil, fmt.Errorf("init repository err: %w", err)
	}

	zap.L().Info("repository inited")
	srv := service.New(repo.Message)
	zap.L().Info("service inited")
	app.handler = handler.New(srv.Message)

	app.server, err = net.Listen("tcp", fmt.Sprintf(":%v", cfg.App.Port))
	if err != nil {
		return nil, fmt.Errorf("init listener err: %w", err)
	}
	zap.L().Info("listener inited")

	return &app, nil
}

func (a *App) Start() error {
	zap.L().Info("start accepting tcp connections...")

	go func() {
		<-a.ctx.Done()
		a.server.Close()
	}()

	for {
		select {
		case <-a.ctx.Done():
			return nil
		default:
			conn, err := a.server.Accept()
			if err != nil {
				zap.L().Error("accept conn err", zap.Error(err))
				continue
			}

			go a.handler.HandleConn(a.ctx, conn)
		}
	}
}

func (a *App) Stop() error {
	zap.L().Info("stopping application...")

	err := a.server.Close()
	if err != nil {
		return err
	}

	err = a.connPool.Close()
	if err != nil {
		return err
	}

	zap.L().Info("application stopped")
	return nil
}
