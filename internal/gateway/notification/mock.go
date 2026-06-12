package notification

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type MockGateway struct {
	log *logrus.Logger
}

func NewMockGateway(log *logrus.Logger) Gateway {
	return &MockGateway{log: log}
}

func (g *MockGateway) Send(ctx context.Context, input Input) (Result, error) {
	g.log.Infof("[MOCK_WA → %s] %s", input.ToPhone, input.Message)
	return Result{
		MessageID: fmt.Sprintf("mock-%s", input.ToPhone),
		Source:    "MOCK_WA",
	}, nil
}

func (g *MockGateway) Source() string {
	return "MOCK_WA"
}
