package main

import (
	"context"

	"github.com/core-go/config"
	"github.com/core-go/log"
	sv "github.com/core-go/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-service/internal/app"
	lg "go-service/pkg/log"
)

func main() {
	var conf app.Root
	er1 := config.Load(&conf, "configs/config")
	if er1 != nil {
		panic(er1)
	}

	e := echo.New()

	log.Initialize(conf.Log)
	logger := lg.NewLogger(conf.MiddleWare, log.InfoFields)
	e.Use(middleware.BodyDump(logger.Log))

	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	fmt.Printf("Request Body: %v\n", string(reqBody))
	// 	fmt.Printf("Response Body: %v\n", string(resBody))
	// 	fmt.Printf("----------------------------------------\n")
	// }))
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "host=${host}, method=${method}, uri=${uri}, status=${status}, error=${error}, message=${message}\n",
	// }))
	e.Use(middleware.Recover())

	er2 := app.Route(e, context.Background(), conf)
	if er2 != nil {
		panic(er2)
	}
	e.Logger.Fatal(e.Start(sv.Addr(conf.Server.Port)))
}
