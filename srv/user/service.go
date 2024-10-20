package user

import (
	"context"
	"errors"
	"sync"

	"git.pmx.cn/hci/microservice-app/proto/user"

)

var (
	ErrUserNotFound = errors.New("user not found")
)

var (
	mem map[int64]*UserInfo
	mu  sync.RWMutex
)

type UserInfo struct {
	UserID  int64
	Name    string
	Company string
	Title   string
}

// NewUserService returns a naive, stateless implementation of Profile Service.
func NewUserService() user.UserServiceServer {
	return service{}
}

type service struct{}

func (s service) GetProfile(c context.Context, req *user.GetProfileRequest) (*user.GetProfileResponse, error) {
	userID := req.GetUserId()
	mu.RLock()
	defer mu.RUnlock()
	if ui, ok := mem[userID]; ok {
		resp := &user.GetProfileResponse{}
		resp.UserId = userID
		resp.Name = ui.Name
		resp.Company = ui.Company
		resp.Title = ui.Title
		//resp.Feeds =
		return resp, nil
	}
	return nil, ErrUserNotFound
}
