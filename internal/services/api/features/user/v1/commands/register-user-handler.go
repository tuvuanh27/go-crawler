package commands

import (
	"context"
	"encoding/json"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/mapper"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	"github.com/tuvuanh27/go-crawler/internal/pkg/repository/interfaces"
	"github.com/tuvuanh27/go-crawler/internal/pkg/utils"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/dtos"
)

type RegisterUserHandler struct {
	log            logger.ILogger
	userRepository interfaces.IUserRepository
	ctx            context.Context
}

func NewRegisterUserHandler(log logger.ILogger, userRepository interfaces.IUserRepository, ctx context.Context) *RegisterUserHandler {
	return &RegisterUserHandler{log: log, userRepository: userRepository, ctx: ctx}
}

func (c *RegisterUserHandler) Handle(ctx context.Context, command *RegisterUser) (*dtos.RegisterUserResponseDto, error) {

	password, err := utils.HashPassword(command.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:     command.Email,
		Password:  password,
		UserName:  command.UserName,
		LastName:  command.LastName,
		FirstName: command.FirstName,
		CreatedAt: command.CreatedAt,
	}

	c.log.Debug("RegisterUser", user)

	registeredUser, err := c.userRepository.RegisterUser(c.ctx, user)
	if err != nil {
		return nil, err
	}

	response, err := mapper.Map[*dtos.RegisterUserResponseDto, *model.User](registeredUser)
	if err != nil {
		return nil, err
	}
	bytes, _ := json.Marshal(response)

	c.log.Info("RegisterUserResponseDto", string(bytes))

	return response, nil
}
