package repository

import (
	"context"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	mongodriver "github.com/tuvuanh27/go-crawler/internal/pkg/mongo-driver"
	"github.com/tuvuanh27/go-crawler/internal/pkg/repository/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	genericRepository *mongodriver.GenericRepository[model.User]
	log               logger.ILogger
	collection        *mongo.Collection
}

func NewUserRepository(db *mongo.Database, log logger.ILogger) interfaces.IUserRepository {
	return &UserRepository{
		genericRepository: mongodriver.NewGenericRepository[model.User](db, mongodriver.UserCollection),
		log:               log,
		collection:        db.Collection(mongodriver.UserCollection),
	}
}

func (u *UserRepository) RegisterUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := u.genericRepository.Add(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	users, err := u.genericRepository.GetAll(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	return users, nil
}
