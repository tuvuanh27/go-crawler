package echomiddleware

import (
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		id := c.Response().Header().Get(echo.HeaderXRequestID)
		if id == "" {
			id = uuid.NewV4().String()
		}

		logWithFields := logger.InitLogger(log.Fields{
			"request_id": id,
		})

		// set logger to context
		c.Set("logger", logWithFields)

		return next(c)
	}
}
