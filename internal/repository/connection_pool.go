package repository

import (
	"os"

	"github.com/Zrossiz/message-broker/internal/config"
)

type ConnectionPool struct {
	File *os.File
}

func NewConnectionPool(cfg config.Repo) (*ConnectionPool, error) {
	var pool ConnectionPool

	file, err := os.OpenFile(cfg.Filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	pool.File = file

	return &pool, nil
}

func (r *ConnectionPool) Close() error {
	err := r.File.Close()
	if err != nil {
		return err
	}

	return nil
}
