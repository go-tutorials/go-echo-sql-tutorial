package app

import (
	"context"

	"github.com/core-go/health"
	s "github.com/core-go/health/sql"
	"github.com/core-go/sql"
	_ "github.com/lib/pq"

	"go-service/internal/handlers"
	"go-service/internal/services"
)

const (
	CreateTable = `
create table if not exists users (
  id varchar(40) not null,
  username varchar(120),
  email varchar(120),
  phone varchar(45),
  date_of_birth date,
  primary key (id)
)`
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	UserHandler   *handlers.UserHandler
}

func NewApp(context context.Context, root Root) (*ApplicationContext, error) {
	db, err := sql.OpenByConfig(root.Sql)
	if err != nil {
		return nil, err
	}

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	sqlChecker := s.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
