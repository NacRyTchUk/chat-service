package datastruct

const UsersChatsTableName = "users_chats"

type UserChats struct {
	UserId int64 `db:"user_id"`
	ChatId int64 `db:"chat_id"`
}
