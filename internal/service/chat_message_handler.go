package service

import (
	"chat-service/internal/dto"
	"log"
)

func (c chatService) ChatHandler(message dto.FormerMessage) error {
	log.Println("chat handler call: ", message)
	msg, err := c.dao.NewMessageQuery().NewTx(message.ChatId, message.SenderName, message.Text)
	if err != nil {
		return err
	}
	c.broker.Publish(dto.Message{
		Id:         msg.Id,
		ChatId:     msg.ChatId,
		SenderName: message.SenderName,
		Text:       msg.Text,
		Timestamp:  msg.CreatedAt.Unix(),
	})
	return nil
}
