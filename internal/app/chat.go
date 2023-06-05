package app

import (
	"chat-service/internal/dto"
	"chat-service/pkg"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	READ_MODE  = int64(1)
	WRITE_MODE = int64(2)
)

type ChatConnection struct {
	ChatMode   int64
	Name       string
	Connection *websocket.Conn
}

func (server *ChatServiceServer) Chat(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	mt, message, err := c.ReadMessage()
	if err != nil || mt == websocket.CloseMessage {
		return
	}

	handshake, err := pkg.Deserialize[HandshakeRequest](message)
	if err != nil {
		return
	}
	log.Println("new handshake: ", handshake)

	chatConn := ChatConnection{
		ChatMode:   handshake.ChatMode,
		Name:       handshake.Name,
		Connection: c,
	}

	log.Println(server.chatService.NewUser(handshake.Name))

	server.chatConnections.NewConnection(chatConn)
	defer server.chatConnections.CloseConnection(chatConn)

	for {
		mt, message, err := c.ReadMessage()
		log.Println("new message")

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		obj, err := pkg.Deserialize[MessageRequest](message)
		if err != nil {
			return
		}
		log.Println("new request: ", obj)

		go func() {
			err := server.chatService.ChatHandler(dto.FormerMessage{
				ChatId:     obj.ChatId,
				SenderName: obj.Name,
				Text:       obj.Text,
			})
			if err != nil {
				return
			}

		}()
	}
}

func (server *ChatServiceServer) BroadcastMessage(message MessageResponse) {
	users, err := server.chatService.GetUsers(message.ChatId)
	if err != nil {
		return
	}

	contains := func(name string) bool {
		for _, user := range users {
			if user.Name == name {
				return true
			}
		}
		return false
	}

	for _, conn := range server.chatConnections.GetConnections() {
		if conn.ChatMode&READ_MODE > 0 && contains(conn.Name) {
			bmsg, err := pkg.Serialize[MessageResponse](message)
			if err != nil {
				return
			}
			_ = conn.Connection.WriteMessage(websocket.BinaryMessage, bmsg)
		}
	}
}
