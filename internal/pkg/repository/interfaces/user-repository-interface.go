package interfaces

import (
	"context"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
)

type IUserRepository interface {
	RegisterUser(ctx context.Context, user *model.User) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
}
