package main

import (
	"chat-service/internal/app"
	"chat-service/internal/broker"
	"chat-service/internal/repository"
	"chat-service/internal/service"
	pb "chat-service/pkg/gen/go/api/service/v1"
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Prepare config file
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	// Connect to db
	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	// Connect to a broker
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", viper.GetString("nats-host"), viper.GetInt("nats-port")))
	if err != nil {
		log.Fatalln(err)
	}
	bkr := broker.NewMessageBroker(nc)

	// Register services
	dao := repository.NewDao()
	chatService := service.NewChatService(dao, bkr)
	chatServer := app.NewChatServiceServer(chatService)

	subscribe, err := nc.Subscribe("new-messages", chatServer.GetMessageHandler())
	if err != nil {
		return
	}
	defer subscribe.Unsubscribe()

	go func() {
		// Open http gateway
		mux := runtime.NewServeMux()

		err = pb.RegisterChatServiceHandlerServer(context.Background(), mux, chatServer)
		if err != nil {
			log.Fatalln("cannot register this service")
		}

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", viper.GetInt("chat-rest-port")),
			Handler: mux,
		}
		log.Fatalln(srv.ListenAndServe())
	}()

	log.Println("Start listen")
	log.Fatalln(chatServer.Listen())
}
