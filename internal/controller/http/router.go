package http

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/myrat012/testing-websocket-golang/internal/service"
	"github.com/myrat012/testing-websocket-golang/internal/usecase"
	"github.com/myrat012/testing-websocket-golang/internal/usecase/repo"
	"github.com/myrat012/testing-websocket-golang/pkg/postgres"
)

func NewRouter() *mux.Router {
	ctx := context.Background()
	dbConn := postgres.Connect(ctx)
	rp := repo.NewInMemoryMessageRepository(dbConn, ctx)
	uc := usecase.NewMessageUseCase(rp)
	wsService := service.NewWebSocketService(uc, dbConn, ctx)

	r := mux.NewRouter()
	r.HandleFunc("/ws", wsService.WebSocketHandler)
	return r
}
