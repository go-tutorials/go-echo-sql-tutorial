package main

import (
	"context"

	"github.com/core-go/config"
	sv "github.com/core-go/core"
	"github.com/core-go/log"
	"github.com/core-go/log/strings"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-service/internal/app"
	mid "go-service/pkg/log"
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
	e.Logger.Fatal(e.Start(sv.Addr(conf.Server.Port)))
}

func MaskLog(name, s string) string {
	if name == "mobileNo" {
		return strings.Mask(s, 2, 2, "x")
	} else {
		return strings.Mask(s, 0, 5, "x")
	}
}
