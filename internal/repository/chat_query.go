package repository

import (
	"chat-service/internal/datastruct"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/blockloop/scan/v2"
)

type ChatQuery interface {
	Create(name string) (*datastruct.Chat, error)
	GetId(name string) (*datastruct.Chat, error)
	GetName(id int64) (*datastruct.Chat, error)
	GetList() ([]datastruct.Chat, error)
}

type chatQuery struct {
	dao dao
	db  *sql.DB
}

func (d dao) NewChatQuery() ChatQuery {
	return chatQuery{d, db}
}

func (c chatQuery) Create(name string) (*datastruct.Chat, error) {
	//TODO implement me
	panic("implement me")
}

func (c chatQuery) GetId(name string) (*datastruct.Chat, error) {
	qb := pgQb().
		Select("id").
		From(datastruct.ChatTableName).
		Where(squirrel.Eq{"name": name})
	var id int64
	err := qb.QueryRow().Scan(&id)
	if err != nil {
		return nil, err
	}
	return &datastruct.Chat{
		Id:   id,
		Name: name,
	}, nil
}

func (c chatQuery) GetName(id int64) (*datastruct.Chat, error) {
	qb := pgQb().
		Select("name").
		From(datastruct.ChatTableName).
		Where(squirrel.Eq{"id": id})
	var name string
	err := qb.QueryRow().Scan(&name)
	if err != nil {
		return nil, err
	}
	return &datastruct.Chat{
		Id:   id,
		Name: name,
	}, nil
}

func (c chatQuery) GetList() ([]datastruct.Chat, error) {
	var items []datastruct.Chat

	tx, err := c.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get lines
	sb := pgQb().
		Select("*").
		From(datastruct.ChatTableName)
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
