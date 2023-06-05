package datastruct

const ChatTableName = "chats"

type Chat struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}
