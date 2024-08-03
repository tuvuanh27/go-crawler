package configurations

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/tuvuanh27/go-crawler/internal/services/api/docs"
)

func ConfigSwagger(e *echo.Echo) {

	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Identities Service Api"
	docs.SwaggerInfo.Description = "Identities Service Api"
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
