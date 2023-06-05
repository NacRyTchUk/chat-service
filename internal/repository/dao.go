package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"

	sq "github.com/Masterminds/squirrel"
	"github.com/spf13/viper"
)

type DAO interface {
	NewMessageQuery() MessageQuery
	NewChatQuery() ChatQuery
	NewUserQuery() UserQuery
	NewUserChatsQuery() UserChatsQuery
}

type dao struct {
}

func NewDao() DAO {
	return &dao{}
}

var db *sql.DB

func pgQb() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(db)
}

func NewDB() (*sql.DB, error) {
	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		viper.Get("database-user").(string),
		viper.Get("database-password").(string),
		viper.Get("database-host").(string),
		viper.Get("database-port").(int),
		viper.Get("database-dbname").(string))
	conn, err := sql.Open("pgx", cs)
	if err != nil {
		return nil, err
	}
	db = conn
	return conn, nil
}
