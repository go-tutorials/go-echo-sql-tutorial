package main

import (
	"context"

	"github.com/core-go/config"
	"github.com/core-go/core"
	"github.com/core-go/core/log"
	mid "github.com/core-go/core/middleware/echo"
	"github.com/core-go/core/strings"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-service/internal/app"
)

func main() {
	var conf app.Config
	err := config.Load(&conf, "configs/config")
	if err != nil {
		panic(err)
	}

	e := echo.New()
	log.Initialize(conf.Log)
	echoLogger := mid.NewEchoLogger(conf.MiddleWare, log.InfoFields, MaskLog)

	e.Use(echoLogger.BuildContextWithMask)
	e.Use(echoLogger.Logger)
	e.Use(middleware.Recover())

	err = app.Route(e, context.Background(), conf)
	if err != nil {
		panic(err)
	}
	e.Logger.Fatal(e.Start(core.Addr(conf.Server.Port)))
}

func MaskLog(name, s string) string {
	if name == "mobileNo" {
		return strings.Mask(s, 2, 2, "x")
	} else {
		return strings.Mask(s, 0, 5, "x")
	}
}
