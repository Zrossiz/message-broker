package repository

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/Zrossiz/message-broker/internal/model"
)

type MessageRepo struct {
	filepath string
}

func NewMessage(filepath string) *MessageRepo {
	return &MessageRepo{
		filepath: filepath,
	}
}

func (r *MessageRepo) Append(msg model.Message) error {
	file, err := os.OpenFile(r.filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = file.Write(append(jsonMsg, '\n'))
	if err != nil {
		return err
	}

	return nil
}

func (r *MessageRepo) Get(offset uint64) (*model.Message, error) {
	file, err := os.OpenFile(r.filepath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	line, err := r.findLine(file, offset)
	if err != nil {
		return nil, err
	}

	var msg model.Message
	err = json.Unmarshal([]byte(line), &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
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
