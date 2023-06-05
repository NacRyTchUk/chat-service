package dto

type Message struct {
	Id         int64
	ChatId     int64
	SenderName string
	Text       string
	Timestamp  int64
}

type FormerMessage struct {
	ChatId     int64
	SenderName string
	Text       string
}

type Chat struct {
	Id   int64
	Name string
}

type User struct {
	Id   int64
	Name string
}
