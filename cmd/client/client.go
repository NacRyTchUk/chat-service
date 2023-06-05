package main

import (
	"chat-service/cmd/client/model"
	"chat-service/internal/app"
	websocket2 "chat-service/pkg"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Chat struct {
	Id   string
	Name string
}

type Message struct {
	Id         string
	ChatId     string
	SenderName string
	Text       string
	Timestamp  string
}

type ListResponse struct {
	Chats []Chat
}
type JoinResponse struct {
	Chat     Chat
	Messages []Message
}

var (
	restPort = 0
)

func ChatJoin(name, chatName string, client *model.Client) {
	chat, err := client.GetChat(chatName)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return
	}
	chatId := chat.Id

	requestURL := fmt.Sprintf("http://localhost:%d/join?name=%s&chatId=%d", restPort, name, chatId)
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	response, err := websocket2.Deserialize[JoinResponse](resBody)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("====History of [%s]====\n", response.Chat.Name)
	for _, msg := range response.Messages {
		id, _ := strconv.Atoi(msg.Id)
		times, _ := strconv.Atoi(msg.Timestamp)
		m := model.Message{
			Id:        int64(id),
			Name:      msg.SenderName,
			Text:      msg.Text,
			Timestamp: int64(times),
		}
		t := time.Unix(m.Timestamp, 0)
		fmt.Printf("~[%d:%d:%d] %s: %s\n", t.Hour(), t.Minute(), t.Second(), msg.SenderName, msg.Text)
		_ = client.NewMessage(app.MessageResponse{
			Id:         m.Id,
			ChatId:     chatId,
			SenderName: m.Name,
			Text:       m.Text,
			Timestamp:  m.Timestamp,
		})
	}
	fmt.Println("======================")
}

func ChatList(client *model.Client) {
	requestURL := fmt.Sprintf("http://localhost:%d/list", restPort)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	response, err := websocket2.Deserialize[ListResponse](resBody)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Println("====Chats available===")
	for _, chat := range response.Chats {
		fmt.Printf("-[%s]\n", chat.Name)
		id, _ := strconv.Atoi(chat.Id)
		_ = client.NewChat(int64(id), chat.Name)
	}
	fmt.Println("======================")
	fmt.Println()
}

func main() {
	// Prepare config file
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	restPort = viper.GetInt("chat-rest-port")

	var name, mode string
	fmt.Print("Enter your name: ")
	fmt.Scanln(&name)
	fmt.Print("Chat mode: ")
	fmt.Scanln(&mode)

	msgRec := make(chan app.MessageResponse, 1000)
	msgSend := make(chan app.MessageRequest, 100)

	client := model.NewClient(name, msgRec, msgSend)

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:" + viper.GetString("chat-primary-port"), Path: "/chat"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	var chatmode int64
	if mode == "w" {
		chatmode = app.WRITE_MODE
	} else {
		chatmode = app.READ_MODE
	}
	bshake, err := websocket2.Serialize[app.HandshakeRequest](app.HandshakeRequest{
		ChatMode: chatmode,
		Name:     name,
	})
	if err != nil {
		return
	}
	err = c.WriteMessage(websocket.BinaryMessage, bshake)
	if err != nil {
		log.Println("write:", err)
		return
	}

	ChatList(client)

	done := make(chan struct{})

	if chatmode == app.WRITE_MODE {

		var chatName string
		fmt.Print("Enter chat name you want to join: ")
		fmt.Scanln(&chatName)
		ChatJoin(name, chatName, client)

		for {
			select {
			case <-done:
				return
			case t := <-msgSend:
				bmsg, _ := json.Marshal(t)
				err := c.WriteMessage(websocket.BinaryMessage, bmsg)
				if err != nil {
					log.Println("write:", err)
					return
				}
			case <-interrupt:
				log.Println("interrupt")

				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			default:
				var msg string
				fmt.Print(name + "->[" + chatName + "]: ")
				fmt.Scanln(&msg)
				chat, err := client.GetChat(chatName)
				if err != nil {
					return
				}
				client.SendMessage(chat.Id, msg)
			}
		}
	} else {
		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				var msg app.MessageResponse
				_ = json.Unmarshal(message, &msg)
				msgRec <- msg
			}
		}()

		go client.Listen()

		for {
			select {
			case <-done:
				return
			case <-interrupt:
				log.Println("interrupt")

				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}

}
