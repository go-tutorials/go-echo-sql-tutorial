package app

import (
	"context"
	"github.com/core-go/health/echo"
	s "github.com/core-go/health/sql"
	"github.com/core-go/sql"
	_ "github.com/lib/pq"

	"go-service/internal/handler"
	"go-service/internal/service"
)

type ApplicationContext struct {
	Health  *echo.Handler
	Handler *handler.UserHandler
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	db, err := sql.OpenByConfig(cfg.Sql)
	if err != nil {
		return nil, err
	}

	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	sqlChecker := s.NewHealthChecker(db)
	healthHandler := echo.NewHandler(sqlChecker)

	return &ApplicationContext{
		Health:  healthHandler,
		Handler: userHandler,
	}, nil
}
