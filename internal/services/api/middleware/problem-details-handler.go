package middlewares

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/tuvuanh27/go-crawler/internal/pkg/utils"
)

func ProblemDetailsHandler(error error, c echo.Context) {
	if !c.Response().Committed {
		if _, err := utils.ResolveProblemDetails(c.Response(), c.Request(), error); err != nil {
			log.Error(err)
		}
	}
}
