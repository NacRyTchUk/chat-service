package repository

import (
	"chat-service/internal/datastruct"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
)

type UserQuery interface {
	Create(name string) (*datastruct.User, error)
	GetId(name string) (*datastruct.User, error)
	GetName(id int64) (*datastruct.User, error)
}

type userQuery struct {
	dao dao
	db  *sql.DB
}

func (d dao) NewUserQuery() UserQuery {
	return userQuery{d, db}
}

func (u userQuery) Create(name string) (*datastruct.User, error) {
	var item datastruct.User

	tx, err := u.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = u.dao.NewUserQuery().GetId(name)
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Insert new line
	ib := pgQb().
		Insert(datastruct.UserTableName).
		Columns("name").
		Values(name).
		Suffix("RETURNING *")
	err = ib.QueryRow().Scan(&item.Id, &item.Name)
	if err != nil {
		return nil, fmt.Errorf("error while item insertion: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commit transaction: %w", err)
	}

	return &item, nil
}

func (u userQuery) GetId(name string) (*datastruct.User, error) {
	qb := pgQb().
		Select("id").
		From(datastruct.UserTableName).
		Where(squirrel.Eq{"name": name})
	var id int64
	err := qb.QueryRow().Scan(&id)
	if err != nil {
		return nil, err
	}
	return &datastruct.User{
		Id:   id,
		Name: name,
	}, nil
}

func (u userQuery) GetName(id int64) (*datastruct.User, error) {
	qb := pgQb().
		Select("name").
		From(datastruct.UserTableName).
		Where(squirrel.Eq{"id": id})
	var name string
	err := qb.QueryRow().Scan(&name)
	if err != nil {
		return nil, err
	}
	return &datastruct.User{
		Id:   id,
		Name: name,
	}, nil
}
