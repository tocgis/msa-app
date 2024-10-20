package auth

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"git.pmx.cn/hci/microservice-app/proto/auth"
)

// Storage
var (
	//mem map[int64]map[int64]*auth.AuthRecord
	mu  sync.RWMutex
)

func init() {
	//mem = make(map[int64]map[int64]*feed.FeedRecord)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

// NewAuthService returns a naive, stateless implementation of Feed Service.
func NewAuthService() auth.AuthServiceServer {
	return service{}
}

type service struct{}

func (s service) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {

	var authRes auth.LoginResponse

	authRes.Token = "xxxxxx"
	return &auth.LoginResponse{
		Token: "TOKENxxxxxxx",
	}, nil


	//panic("implement me")
}

func (s service) Auth(ctx context.Context, request *auth.AuthRequest) (*auth.AuthResponse, error) {

	var res auth.AuthResponse
	var info auth.TokenInfo
	info.Name = "xxxx"
	info.Exp = 122342342
	info.Phone = "13922332233"
	info.UserId = 1

	res.Code = http.StatusOK
	res.Data = &info
	res.Message = "OK"

	return &res, nil

	//panic("implement me")
}
