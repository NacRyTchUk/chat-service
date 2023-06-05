package repository

import (
	"chat-service/internal/datastruct"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/blockloop/scan/v2"
)

type UserChatsQuery interface {
	Create(userId int64, chatId int64) (*datastruct.UserChats, error)
	GetChats(userId int64) ([]datastruct.Chat, error)
	GetUsers(chatId int64) ([]datastruct.User, error)
}

type userChatsQuery struct {
	dao dao
	db  *sql.DB
}

func (d dao) NewUserChatsQuery() UserChatsQuery {
	return userChatsQuery{d, db}
}

func (u userChatsQuery) Create(userId int64, chatId int64) (*datastruct.UserChats, error) {
	var item datastruct.UserChats

	tx, err := u.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = u.dao.NewUserQuery().GetName(userId)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	_, err = u.dao.NewChatQuery().GetName(chatId)
	if err != nil {
		return nil, fmt.Errorf("chat not found")
	}

	// Insert new line
	ib := pgQb().
		Insert(datastruct.UsersChatsTableName).
		Columns("user_id, chat_id").
		Values(userId, chatId).
		Suffix("RETURNING *")
	err = ib.QueryRow().Scan(&item.UserId, &item.ChatId)
	if err != nil {
		return nil, fmt.Errorf("error while item insertion: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commit transaction: %w", err)
	}

	return &item, nil
}

func (u userChatsQuery) GetChats(userId int64) ([]datastruct.Chat, error) {
	var items []datastruct.Chat

	tx, err := u.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get lines
	sb := pgQb().
		Select("chat_id").
		From(datastruct.UsersChatsTableName).
		Where(squirrel.Eq{"user_id": userId})
	rows, err := sb.Query()
	if err != nil {
		return nil, fmt.Errorf("error while getting items: %w", err)
	}
	var chatsId []int64
	err = scan.Rows(&chatsId, rows)
	if err != nil {
		return nil, fmt.Errorf("error while rows scanning: %w", err)
	}

	for _, chatId := range chatsId {
		chat, err := u.dao.NewChatQuery().GetName(chatId)
		if err != nil {
			return nil, err
		}
		items = append(items, *chat)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commit transaction: %w", err)
	}
	return items, nil
}

func (u userChatsQuery) GetUsers(chatId int64) ([]datastruct.User, error) {
	var items []datastruct.UserChats

	tx, err := u.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get lines
	sb := pgQb().
		Select("*").
		From(datastruct.UsersChatsTableName).
		Where(squirrel.Eq{"chat_id": chatId})
	rows, err := sb.Query()
	if err != nil {
		return nil, fmt.Errorf("error while getting items: %w", err)
	}
	err = scan.Rows(&items, rows)
	if err != nil {
		return nil, fmt.Errorf("error while rows scanning: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commit transaction: %w", err)
	}

	var users []datastruct.User
	for _, userChat := range items {
		user, _ := u.dao.NewUserQuery().GetName(userChat.UserId)
		users = append(users, datastruct.User{
			Id:   user.Id,
			Name: user.Name,
		})
	}
	return users, nil
}
