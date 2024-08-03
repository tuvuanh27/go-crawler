package configurations

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/product/v1/controllers"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/endpoints"
)

func ConfigEndpoints(validator *validator.Validate, echo *echo.Echo, ctx context.Context) {
	endpoints.MapRoute(validator, echo, ctx)
	controllers.MapRoute(validator, echo, ctx)
}
