package service

import "chat-service/internal/dto"

func (c chatService) ChatList() (list []dto.Chat, err error) {
	chats, err := c.dao.NewChatQuery().GetList()
	if err != nil {
		return nil, err
	}
	for _, chat := range chats {
		list = append(list, dto.Chat(chat))
	}
	return
}
