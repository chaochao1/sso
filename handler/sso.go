package handler

import (
	"context"
	"github.com/micro/go-log"
	sso "github.com/cicdi-go/sso/proto/sso"
	"github.com/cicdi-go/sso/src/models"
	"errors"
	"github.com/cicdi-go/sso/src/utils"
)

type Sso struct{}

// Call is a single request handler called via client.Auth or the generated client code
func (e *Sso) Token(ctx context.Context, req *sso.AuthRequest, rsp *sso.AuthResponse) error {
	u := models.User{
		Username: req.Username,
	}
	if u.Verify(req.Password) {
		err := errors.New("账号或者密码错误")
		return err
	}
	token, expire, err := u.GenerateToken()
	if err != nil {
		return err
	}
	rsp.Expire = expire.Unix()
	rsp.Token = token
	return nil
}

// Call is a single request handler called via client.Register or the generated client code
func (e *Sso) Register(ctx context.Context, req *sso.RegisterRequest, rsp *sso.RegisterResponse) error {
	if !utils.CaptchaVerify(req.CaptchaId, req.Verify) {
		return errors.New("验证码验证失败")
	}
	if req.Username == "" {
		return errors.New("用户名不能为空")
	}
	user := models.User{
		Username: req.Username,
		Status:   1,
	}
	user.SetPassword(req.Password)
	if err := user.Insert(); err != nil {
		return err
	}
	rsp.Username = user.Username
	rsp.Status = int64(user.Status)
	return nil
}

func (e * Sso) Captcha(ctx context.Context, req *sso.CaptchaRequest, rsp *sso.CaptchaResponse) error {
	id, data, err := utils.CaptchaGenerate(req.Type, req.Length)
	if err != nil {
		return err
	}
	rsp.Id = id
	rsp.Data = data
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Sso) Stream(ctx context.Context, req *sso.StreamingRequest, stream sso.Sso_StreamStream) error {
	log.Logf("Received Sso.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&sso.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Sso) PingPong(ctx context.Context, stream sso.Sso_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&sso.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
