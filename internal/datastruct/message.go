package datastruct

import "time"

const MessageTableName = "messages"

type Message struct {
	Id        int64     `db:"id"`
	ChatId    int64     `db:"chat_id"`
	SenderId  int64     `db:"sender_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}
