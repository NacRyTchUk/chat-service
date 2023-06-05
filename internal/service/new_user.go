package service

import "chat-service/internal/dto"

func (c chatService) NewUser(name string) (dto.User, error) {
	user, err := c.dao.NewUserQuery().Create(name)
	if err != nil {
		return dto.User{}, err
	}
	return dto.User(*user), nil
}
