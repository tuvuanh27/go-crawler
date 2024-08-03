package endpoints

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/commands"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/dtos"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/queries"
	"net/http"
)

func MapRoute(validator *validator.Validate, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/users")
	group.POST("", createUser(validator, ctx))
	group.GET("", getUsers(validator, ctx))
}

// RegisterUser
// @Tags Users
// @Summary Register user
// @Description Create new user
// @Accept json
// @Produce json
// @Param RegisterUserRequestDto body dtos.RegisterUserRequestDto true "User data"
// @Success 201 {object} dtos.RegisterUserResponseDto
// @Security ApiKeyAuth
// @Router /api/v1/users [post]
func createUser(validator *validator.Validate, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := c.Get("logger").(logger.ILogger)
		request := &dtos.RegisterUserRequestDto{}

		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[registerUserEndpoint_handler.Bind] error in the binding request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := commands.NewRegisterUser(request.FirstName, request.LastName, request.UserName, request.Email, request.Password)

		if err := validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[registerUserEndpoint_handler.StructCtx] command validation failed")
			log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result, err := mediatr.Send[*commands.RegisterUser, *dtos.RegisterUserResponseDto](ctx, command)

		if err != nil {
			log.Errorf("(RegisterUser.Handle) id: {%d}, err: {%v}", result.ID, err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		log.Infof("(user registered) id: {%d}", result.ID)
		return c.JSON(http.StatusCreated, result)
	}
}

func getUsers(validator *validator.Validate, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		query := queries.NewGetUsersQuery()
		log := c.Get("logger").(logger.ILogger)

		result, err := mediatr.Send[*queries.GetUsersQuery, []*dtos.GetUsersResponseDto](ctx, query)
		if err != nil {
			log.Errorf("(GetUsers.Handle) err: {%v}", err)
		}

		return c.JSON(http.StatusOK, result)

	}
}
