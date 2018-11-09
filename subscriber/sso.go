package subscriber

import (
	"context"
	"github.com/micro/go-log"

	sso "github.com/cicdi-go/sso/proto/sso"
)

type Sso struct{}

func (e *Sso) Handle(ctx context.Context, msg *sso.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *sso.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
