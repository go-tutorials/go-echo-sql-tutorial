package app

import (
	"context"

	"github.com/labstack/echo/v4"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

func Route(e *echo.Echo, ctx context.Context, root Root) error {
	app, err := NewApp(ctx, root)
	if err != nil {
		return err
	}

	e.GET("/health", app.HealthHandler.Check)

	userPath := "/users"
	e.GET(userPath, app.UserHandler.GetAll)
	e.GET(userPath+"/:id", app.UserHandler.Load)
	e.POST(userPath, app.UserHandler.Insert)
	e.PUT(userPath+"/:id", app.UserHandler.Update)
	e.PATCH(userPath+"/:id", app.UserHandler.Patch)
	e.DELETE(userPath+"/:id", app.UserHandler.Delete)

	return nil
}
