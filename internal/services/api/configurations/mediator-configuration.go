package configurations

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/rabbitmq"
	"github.com/tuvuanh27/go-crawler/internal/pkg/repository/interfaces"
	productcommands "github.com/tuvuanh27/go-crawler/internal/services/api/features/product/v1/commands/crawl-aliexpress-product"
	productdtos "github.com/tuvuanh27/go-crawler/internal/services/api/features/product/v1/dtos"
	productqueries "github.com/tuvuanh27/go-crawler/internal/services/api/features/product/v1/queries/get-products"
	usercommands "github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/commands"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/dtos"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/queries"
	"github.com/tuvuanh27/go-crawler/internal/services/api/service"
)

func ConfigUsersMediator(log logger.ILogger, userRepository interfaces.IUserRepository, ctx context.Context, publisher rabbitmq.IPublisher) error {
	if err := mediatr.RegisterRequestHandler[*usercommands.RegisterUser, *dtos.RegisterUserResponseDto](usercommands.NewRegisterUserHandler(log, userRepository, ctx)); err != nil {
		return err
	}

	if err := mediatr.RegisterRequestHandler[*queries.GetUsersQuery, []*dtos.GetUsersResponseDto](queries.NewGetUsersHandler(log, userRepository, ctx, publisher)); err != nil {
		return err
	}

	return nil
}

func ConfigProductMediator(log logger.ILogger, productRepository interfaces.IProductRepository, ctx context.Context, queueService service.IQueueService) error {
	if err := mediatr.RegisterRequestHandler[*productcommands.CrawlAliexpressProduct, *string](productcommands.NewCrawlAliexpressProductHandler(log, ctx, queueService)); err != nil {
		return err
	}

	if err := mediatr.RegisterRequestHandler[*productqueries.GetProductsQuery, *productdtos.GetPaginationByTypeResponseDto](productqueries.NewGetProductsHandler(log, productRepository, ctx)); err != nil {
		return err
	}

	return nil
}
