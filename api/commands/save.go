package commands

import (
	"errors"

	"github.com/gustavosvalentim/evilcode/api/ports"
)

type SaveCommand struct {
	Filename string
	Content  string
}

func NewSaveCommand(filename, content string) (*SaveCommand, error) {
	if len(filename) == 0 {
		return nil, errors.New("Filename can't be empty")
	}
	return &SaveCommand{
		Filename: filename,
		Content:  content,
	}, nil
}

type SaveCommandHandler struct {
	storage ports.FileStorageRepo
}

func NewSaveCommandHandler(storage ports.FileStorageRepo) *SaveCommandHandler {
	return &SaveCommandHandler{
		storage: storage,
	}
}

func (handler *SaveCommandHandler) Handle(command *SaveCommand) error {
	if err := handler.storage.Save(command.Filename, command.Content); err != nil {
		return err
	}
	return nil
}
