package model

import (
	"chat-service/internal/app"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Message struct {
	Id        int64
	Name      string
	Text      string
	Timestamp int64
}

type Chat struct {
	Id       int64
	Name     string
	Messages []Message
	SeenId   int64
}

func (cl *Client) NewChat(id int64, name string) error {
	cl.savedChats[id] = Chat{
		Id:       id,
		Name:     name,
		Messages: nil,
		SeenId:   0,
	}
	return nil
}

func (cl *Client) GetChat(name string) (*Chat, error) {
	for _, v := range cl.savedChats {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, errors.New("not found")
}

func (cl *Client) GetChatById(id int64) (*Chat, error) {
	for _, v := range cl.savedChats {
		if v.Id == id {
			return &v, nil
		}
	}
	return nil, errors.New("not found")
}

func (cl *Client) NewMessage(message app.MessageResponse) error {
	cl.mut.Lock()
	defer cl.mut.Unlock()

	chat := cl.savedChats[message.ChatId]
	chat.Messages = append(chat.Messages, Message{
		Id:        message.Id,
		Name:      message.SenderName,
		Text:      message.Text,
		Timestamp: message.Timestamp,
	})

	return nil
}

type ReceiveMessage struct {
	ChatId int64
	Msg    Message
}

type SendMessage struct {
	ChatId int64
	Name   string
	Text   string
}

type Client struct {
	name       string
	savedChats map[int64]Chat
	mut        sync.RWMutex
	rec        <-chan app.MessageResponse
	snd        chan<- app.MessageRequest
}

var client Client

func NewClient(name string,
	rec <-chan app.MessageResponse,
	snd chan<- app.MessageRequest,
) *Client {
	client = Client{
		name:       name,
		savedChats: nil,
		rec:        rec,
		snd:        snd,
	}
	client.savedChats = make(map[int64]Chat)
	client.Listen()
	return &client
}

func (cl *Client) Listen() {
	go func() {
		for {
			if err := cl.ReceiveMessage(); err != nil {
				log.Fatal(fmt.Errorf("receive message: %w", err))
			}
		}
	}()
}

func (cl *Client) SendMessage(chatId int64, text string) {
	cl.snd <- app.MessageRequest{
		ChatId: chatId,
		Name:   cl.name,
		Text:   text,
	}
}

func (cl *Client) ReceiveMessage() error {
	msg, ok := <-cl.rec
	if !ok {
		return fmt.Errorf("chan closed")
	}
	_ = cl.NewMessage(msg)
	t := time.Unix(msg.Timestamp, 0)
	chat, err := cl.GetChatById(msg.ChatId)
	if err != nil {
		return err
	}
	fmt.Printf("[%d:%d:%d] [%s] %s: %s\n", t.Hour(), t.Minute(), t.Second(), chat.Name, msg.SenderName, msg.Text)
	return nil
}
