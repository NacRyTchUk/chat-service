package service

import "chat-service/internal/dto"

func (c chatService) GetUsers(chatId int64) ([]dto.User, error) {
	users, err := c.dao.NewUserChatsQuery().GetUsers(chatId)
	if err != nil {
		return nil, err
	}
	var chatUsers []dto.User
	for _, user := range users {
		chatUsers = append(chatUsers, dto.User{
			Id:   user.Id,
			Name: user.Name,
		})
	}
	return chatUsers, nil
}
