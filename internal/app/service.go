package app

import (
	"chat-service/internal/dto"
	"chat-service/internal/service"
	"chat-service/pkg"
	pb "chat-service/pkg/gen/go/api/service/v1"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

type ChatServiceServer struct {
	pb.UnimplementedChatServiceServer
	chatService     service.ChatService
	chatConnections ConnectionPool[ChatConnection, bool]
}

func NewChatServiceServer(chatService service.ChatService) *ChatServiceServer {
	return &ChatServiceServer{chatService: chatService}
}

func (server *ChatServiceServer) GetMessageHandler() nats.MsgHandler {
	return func(msg *nats.Msg) {
		log.Println("hand a message")
		message, err := pkg.Deserialize[dto.Message](msg.Data)
		if err != nil {
			return
		}
		server.BroadcastMessage(MessageResponse(message))
	}
}

func (server *ChatServiceServer) Listen() error {
	http.HandleFunc("/chat", server.Chat)
	//http.HandleFunc("/join", server.ChatJoin)
	//http.HandleFunc("/list", server.ChatList)
	if err := http.ListenAndServe(":"+viper.GetString("chat-primary-port"), nil); err != nil {
		log.Println("try to connect to backup server")
		return http.ListenAndServe(":"+viper.GetString("chat-backup-port"), nil)
	}
	return nil
}
