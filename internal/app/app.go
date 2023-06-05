package app

import (
	"chat-service/internal/dto"
	"chat-service/pkg"
	"io"
	"log"
	"net/http"
)

type MessageResponse struct {
	Id         int64
	ChatId     int64
	SenderName string
	Text       string
	Timestamp  int64
}

type MessageRequest struct {
	ChatId int64
	Name   string
	Text   string
}

type HandshakeRequest struct {
	ChatMode int64
	Name     string
}

type HandshakeResponse struct {
	//...
}

type ChatListRequest struct {
	//...
}

type ChatListResponse struct {
	List []dto.Chat
}

type ChatJoinRequest struct {
	Name   string
	ChatId string
}

type ChatJoinResponse struct {
	Chat         dto.Chat
	LastMessages []dto.Message
}

func GetRequestBody[T any](r *http.Request) (body T, err error) {
	log.Println(r.Body)
	//bodyReader, err := r.GetBody()
	//if err != nil {
	//	return
	//}
	bodyRaw, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	return pkg.Deserialize[T](bodyRaw)
}

func SetResponseBody[T any](w http.ResponseWriter, body T) (err error) {
	bodyRaw, err := pkg.Serialize[T](body)
	if err != nil {
		return err
	}
	_, err = w.Write(bodyRaw)
	return
}
