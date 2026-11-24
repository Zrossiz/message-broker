package repository

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type MessageRepo struct {
	curOffset uint64
	file      *os.File
}

func NewMessage(file *os.File) *MessageRepo {
	return &MessageRepo{
		file: file,
	}
}

func (r *MessageRepo) Append(msg []byte) error {
	_, err := r.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	_, err = r.file.Write(msg)
	if err != nil {
		return err
	}

	return nil
}

func (r *MessageRepo) Get(offset uint64) ([]byte, error) {
	line, err := r.findLine(r.file, offset)
	if err != nil {
		return nil, err
	}

	var msg []byte
	err = json.Unmarshal([]byte(line), &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *MessageRepo) findLine(file *os.File, offset uint64) (string, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return "", nil
	}

	scanner := bufio.NewScanner(file)
	var line uint64
	for scanner.Scan() {
		if line == offset {
			return scanner.Text(), nil
		}
		line++
	}

	if scanner.Err() != nil {
		return "", err
	}

	return "", io.EOF
}
