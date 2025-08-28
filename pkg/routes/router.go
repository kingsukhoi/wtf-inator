package routes

import (
	"github.com/kingsukhoi/wtf-inator/pkg/conf"
	"github.com/kingsukhoi/wtf-inator/pkg/proxy"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter() (*echo.Echo, error) {
	config := conf.MustGetConfig()

	currProxy, err := proxy.NewWtfProxy(config.Server.Url)
	if err != nil {
		return nil, err
	}

	e := echo.New()
	e.Use(middleware.Recover())

	e.Any("/*", currProxy.RequestHandler)

	return e, nil

}
