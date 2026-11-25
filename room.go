package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type room struct {

	//будет как бы ужерживать всех клиентво в комнате
	clients map[*client]bool

	//канал для всех тех кто хочет попасть в эту комнату??
	join chan *client

	//те кто хочет ливнуть
	leave chan *client

	//канал который удерживает входящие сообщение, что должны быть адресованы другим клиентам
	forward chan []byte
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

// каждая комната отдельно от другой и может длиться пока сервер жив
func (r *room) run() {
	for {
		select {
		//добавляем пользователя в канал
		case client := <-r.join:
			r.clients[client] = true
			//удаляем из канала
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.recieve)
			//напишем всем что чел ливнул
		case msg := <-r.forward:
			for client := range r.clients {
				client.recieve <- msg
			}
		}
	}
}

//обновляем базовое http соединение на ws

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

var rooms = make(map[string]*room)
var mu sync.Mutex

func getRoom(name string) *room {
	mu.Lock()
	defer mu.Unlock()
	if r, ok := rooms[name]; ok {
		return r
	}
	r := newRoom()
	rooms[name] = r
	go r.run()
	return r
}

//handshake

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	roomName := req.URL.Query().Get("room")
	if roomName == "" {
		http.Error(w, "Room name required", http.StatusBadRequest)
		return
	}

	realRoom := getRoom(roomName)

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &client{
		socket:  socket,
		recieve: make(chan []byte, messageBufferSize),
		room:    realRoom,
		name:    fmt.Sprintf("user_%d", rand.Intn(1000)),
	}

	realRoom.join <- client
	defer func() { realRoom.leave <- client }()
	go client.write()
	client.read()
}
