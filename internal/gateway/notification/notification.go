package notification

import "context"

type Input struct {
	ToPhone string
	Message string
}

type Result struct {
	MessageID string
	Source    string
}

type Gateway interface {
	Send(ctx context.Context, input Input) (Result, error)
	Source() string
}
