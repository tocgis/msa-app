package assess

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"git.pmx.cn/hci/microservice-app/proto/assess"
)

var (
	ErrTopicNotFound = errors.New("topic not found")
)

var (
	mem map[int64]*Assess
	mu  sync.RWMutex
)

type Assess struct {
	UserId int64
	Score string
	Phone string
}

// NewAssessService returns a naive, stateless implementation of Topic Service.
func NewAssessService() assess.AssessServiceServer {
	return service{}
}

type service struct{}

func (s service) InitScore(ctx context.Context, request *assess.ScoreRequest) (*assess.OkResponse, error) {

	return &assess.OkResponse{}, nil
}

func (s service) ScoreInfo(ctx context.Context, request *assess.ScoreRequest) (*assess.ScoreResponse, error) {
	fmt.Println(request)
	return &assess.ScoreResponse{
		SocialScore: "850",
		FinancialScore: "350",
	}, nil
}

func (s service) BasicSave(ctx context.Context, request *assess.BasicInfoRequest) (*assess.OkResponse, error) {
	fmt.Println(request)

	return &assess.OkResponse{}, nil
}

func (s service) EducationSave(ctx context.Context, request *assess.EducationRequest) (*assess.OkResponse, error) {
	return &assess.OkResponse{}, nil
}

func (s service) WorkSave(ctx context.Context, request *assess.WorkinfoRequest) (*assess.OkResponse, error) {
	return &assess.OkResponse{}, nil
}

