package backend

import (
	"context"

	"github.com/rancher/rio/controllers/backend/stack"
	"github.com/rancher/rio/types"
)

func Register(ctx context.Context, rContext *types.Context) error {
	stack.Register(ctx, rContext)
	return nil
}
