package app

import (
	v1 "chat-service/pkg/gen/go/api/model/v1"
	pb "chat-service/pkg/gen/go/api/service/v1"
	"context"
	"log"
)

func (server *ChatServiceServer) Join(ctx context.Context, request *pb.JoinRequest) (*pb.JoinResponse, error) {
	chat, messages, err := server.chatService.ChatJoin(request.Name, request.ChatId)
	log.Println("cme ", chat, messages)
	if err != nil {
		return nil, err
	}
	var msgs []*v1.Message
	for _, msg := range messages {
		msgs = append(msgs, &v1.Message{
			Id:         msg.Id,
			ChatId:     msg.ChatId,
			SenderName: msg.SenderName,
			Text:       msg.Text,
			Timestamp:  msg.Timestamp,
		})
	}

	return &pb.JoinResponse{
		Chat: &v1.Chat{
			Id:   chat.Id,
			Name: chat.Name,
		},
		Messages: msgs,
	}, nil
}
