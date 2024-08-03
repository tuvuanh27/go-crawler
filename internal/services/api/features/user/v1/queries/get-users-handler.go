package queries

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/rabbitmq"
	"github.com/tuvuanh27/go-crawler/internal/pkg/repository/interfaces"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/dtos"
	"time"
)

type GetUsersHandler struct {
	log               logger.ILogger
	userRepository    interfaces.IUserRepository
	ctx               context.Context
	rabbitmqPublisher rabbitmq.IPublisher
}

func NewGetUsersHandler(log logger.ILogger, userRepository interfaces.IUserRepository, ctx context.Context, rabbitmqPublisher rabbitmq.IPublisher) *GetUsersHandler {
	return &GetUsersHandler{log: log, userRepository: userRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher}
}

type CrawlAliexpressProductPayload struct {
	ProductID uuid.UUID `json:"product_id"`
	Url       string    `json:"url"`
}

func (c *GetUsersHandler) Handle(ctx context.Context, query *GetUsersQuery) ([]*dtos.GetUsersResponseDto, error) {

	users, err := c.userRepository.GetAll(c.ctx)
	if err != nil {
		return nil, err
	}

	var response []*dtos.GetUsersResponseDto
	for _, user := range users {
		response = append(response, &dtos.GetUsersResponseDto{
			ID:        user.ID,
			Email:     user.Email,
			UserName:  user.UserName,
			LastName:  user.LastName,
			FirstName: user.FirstName,
			CreatedAt: user.CreatedAt,
		})
	}

	response = append(response, &dtos.GetUsersResponseDto{
		ID:        "1",
		Email:     "dsfdsf",
		UserName:  "tuvuanh27",
		LastName:  "Nguyen",
		FirstName: "Tu",
		CreatedAt: time.Now(),
	})

	//var msg []dtos.GetUsersResponseDto
	//for _, v := range response {
	//	msg = append(msg, *v)
	//}
	//c.log.Debug("GetUsers", msg)
	//
	//err = c.rabbitmqPublisher.PublishMessage(msg)
	//if err != nil {
	//	return nil, err
	//}

	testMsg := CrawlAliexpressProductPayload{
		ProductID: uuid.NewV4(),
		Url:       "https://www.aliexpress.com/item/1005001862351987.html",
	}
	arrMsg := []CrawlAliexpressProductPayload{testMsg}

	err = c.rabbitmqPublisher.PublishMessage(arrMsg)
	if err != nil {
		return nil, err
	}

	return response, nil
}
