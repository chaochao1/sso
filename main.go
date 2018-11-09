package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/cicdi-go/sso/subscriber"

	sso "github.com/cicdi-go/sso/proto/sso"
	"github.com/cicdi-go/sso/handler"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.sso"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	sso.RegisterSsoHandler(service.Server(), new(handler.Sso))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.sso", service.Server(), new(subscriber.Sso))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.sso", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
