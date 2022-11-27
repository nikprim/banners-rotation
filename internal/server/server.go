package server

import "context"

type IServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
