package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/myrat012/testing-websocket-golang/internal/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketService struct {
	usecase usecase.MessageUseCase
	dbConn  *pgxpool.Pool
	clients map[*websocket.Conn]bool
	notify  chan []byte
	ctx     context.Context
	mu      sync.RWMutex
}

func NewWebSocketService(uc usecase.MessageUseCase, db *pgxpool.Pool, cx context.Context) *WebSocketService {
	ws := &WebSocketService{
		usecase: uc,
		dbConn:  db,
		clients: make(map[*websocket.Conn]bool),
		notify:  make(chan []byte),
		ctx:     cx,
	}
	go ws.listenForNotifications()
	return ws
}

func (ws *WebSocketService) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Conn err: ", err)
		return
	}
	defer conn.Close()

	ws.mu.Lock()
	ws.clients[conn] = true
	ws.mu.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			ws.mu.Lock()
			delete(ws.clients, conn)
			defer ws.mu.Unlock()
			return
		}
	}
}

func (ws *WebSocketService) listenForNotifications() {
	conn, err := ws.dbConn.Acquire(ws.ctx)
	if err != nil {
		log.Fatalf("Failed to acquire connection: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ws.ctx, "LISTEN new_record")
	if err != nil {
		log.Fatalf("Failed to listen to new_record channel: %v", err)
	}

	for {
		notification, err := conn.Conn().WaitForNotification(ws.ctx)
		if err != nil {
			log.Printf("Failed to wait for notification: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		ws.notifyClients([]byte(notification.Payload))
	}
}

func (ws *WebSocketService) notifyClients(message []byte) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	for client := range ws.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			client.Close()
			ws.mu.Lock()
			delete(ws.clients, client)
			ws.mu.Unlock()
		}
	}
}
