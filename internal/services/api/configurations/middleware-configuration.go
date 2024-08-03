package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	echomiddleware "github.com/tuvuanh27/go-crawler/internal/pkg/http/echo/middleware"
	"github.com/tuvuanh27/go-crawler/internal/pkg/otel"
	otelmiddleware "github.com/tuvuanh27/go-crawler/internal/pkg/otel/middlewares"
	"github.com/tuvuanh27/go-crawler/internal/services/api/constants"
	middlewares "github.com/tuvuanh27/go-crawler/internal/services/api/middleware"
	"strings"
	"time"
)

func ConfigMiddlewares(e *echo.Echo, jaegerCfg *otel.JaegerConfig) {

	e.HideBanner = false

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "custom timeout error message returns to client",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Info(c.Path(), "timeout error")
		},
		Timeout: 30 * time.Second,
	}))

	e.Use(otelmiddleware.EchoTracerMiddleware(jaegerCfg.ServiceName))
	e.HTTPErrorHandler = middlewares.ProblemDetailsHandler

	e.Use(echomiddleware.CorrelationIdMiddleware)
	e.Use(echomiddleware.RequestIdMiddleware)

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	e.Use(middleware.BodyLimit(constants.BodyLimit))
	e.Use(echomiddleware.LoggerMiddleware)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' https://ae01.alicdn.com; script-src 'self'; style-src 'self';")
			return next(c)
		}
	})

	// serve static reactjs build
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: nil,
		// Root directory from where the static content is served.
		Root: "../../../frontend/dist",
		// Index file for serving a directory.
		// Optional. Default value "index.html".
		Index: "index.html",
		// Enable HTML5 mode by forwarding all not-found requests to root so that
		// SPA (single-page application) can handle the routing.
		HTML5:      true,
		Browse:     false,
		IgnoreBase: false,
		Filesystem: nil,
	}))

	// cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Adjust to restrict to specific origins if needed
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

}
