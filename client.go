package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// один чел что чатится в комнате
type client struct {
	//веб сокет для этого пользователя
	socket *websocket.Conn

	//канал для принятия сообщений от других клиентов
	recieve chan []byte

	room *room

	name string
}

// функция для отправки сообщений
func (c client) read() {
	//нужно закрыть соединение если мы закончили чатиться
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()

		if err != nil {
			return
		}

		outgoing := map[string]string{
			"name":    c.name,
			"message": string(msg),
		}

		jsMessage, err := json.Marshal(outgoing)
		if err != nil {
			fmt.Println("Enconding failed!")
			continue
		}

		c.room.forward <- jsMessage
	}
}

// для принятия сообщений
func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.recieve {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
