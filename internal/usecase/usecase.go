package usecase

import (
	"fmt"

	"github.com/myrat012/testing-websocket-golang/internal/usecase/repo"
)

type MessageUseCase struct {
	rp repo.MessageRepository
}

func NewMessageUseCase(r repo.MessageRepository) MessageUseCase {
	return MessageUseCase{
		rp: r,
	}
}

func (uc *MessageUseCase) ProccessMessage(message string) (string, error) {
	// business logic goes here
	// for example saveing the message to the repository
	err := uc.rp.SaveMessage(message)
	if err != nil {
		fmt.Println("Save message error on usecase: ", err)
		return "", err
	}
	return "Message processed: " + message, nil
}
