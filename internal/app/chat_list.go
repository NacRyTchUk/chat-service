package app

import (
	v1 "chat-service/pkg/gen/go/api/model/v1"
	pb "chat-service/pkg/gen/go/api/service/v1"
	"context"
	"net/http"
)

func (server *ChatServiceServer) List(ctx context.Context, request *pb.ListRequest) (*pb.ListResponse, error) {
	list, err := server.chatService.ChatList()
	if err != nil {
		return nil, err
	}
	var chats []*v1.Chat
	for _, chat := range list {
		chats = append(chats, &v1.Chat{
			Id:   chat.Id,
			Name: chat.Name,
		})
	}
	return &pb.ListResponse{
		Chats: chats,
	}, nil
}

func (server *ChatServiceServer) ChatList(w http.ResponseWriter, r *http.Request) {
	//_, err := GetRequestBody[ChatListRequest](r)
	//if err != nil {
	//	return
	//}
	list, err := server.chatService.ChatList()
	if err != nil {
		return
	}
	err = SetResponseBody[ChatListResponse](w, ChatListResponse{
		List: list,
	})
	if err != nil {
		return
	}
}
