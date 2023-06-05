package repository

import (
	"chat-service/internal/datastruct"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/blockloop/scan/v2"
	"time"
)

type MessageQuery interface {
	NewTx(chat int64, sender, text string) (*datastruct.Message, error)
	Create(chatId, senderId int64, text string) (*datastruct.Message, error)
	GetChatMessages(chatId int64, period time.Duration) ([]datastruct.Message, error)
	GetUserMessages(userId int64, period time.Duration) ([]datastruct.Message, error)
}

type messageQuery struct {
	dao dao
	db  *sql.DB
}

func (d dao) NewMessageQuery() MessageQuery {
	return messageQuery{
		d,
		db,
	}
}

func (m messageQuery) NewTx(chat int64, sender, text string) (*datastruct.Message, error) {
	user, err := m.dao.NewUserQuery().GetId(sender)
	if err != nil {
		return nil, err
	}
	msg, err := m.Create(chat, user.Id, text)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (m messageQuery) Create(chatId, senderId int64, text string) (*datastruct.Message, error) {
	var item datastruct.Message

	tx, err := m.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert new line
	ib := pgQb().
		Insert(datastruct.MessageTableName).
		Columns("chat_id, sender_id, text, created_at").
		Values(chatId, senderId, text, time.Now()).
		Suffix("RETURNING *")
	err = ib.QueryRow().Scan(&item.Id, &item.ChatId, &item.SenderId, &item.Text, &item.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error while item insertion: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commit transaction: %w", err)
	}

	return &item, nil
}

func (m messageQuery) GetChatMessages(chatId int64, period time.Duration) ([]datastruct.Message, error) {
	var items []datastruct.Message

	tx, err := m.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get lines
	sb := pgQb().
		Select("*").
		From(datastruct.MessageTableName).
		Where(squirrel.And{squirrel.Eq{"chat_id": chatId}, squirrel.GtOrEq{"created_at": time.Now().Add(-period)}})
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
	return items, nil
}

func (m messageQuery) GetUserMessages(userId int64, period time.Duration) ([]datastruct.Message, error) {
	//TODO implement me
	panic("implement me")
}
