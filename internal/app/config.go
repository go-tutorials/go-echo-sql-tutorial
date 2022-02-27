package app

import (
	"github.com/core-go/log"
	sv "github.com/core-go/service"
	"github.com/core-go/sql"

	l "go-service/pkg/log"
)

type Config struct {
	Server     sv.ServerConf `mapstructure:"server"`
	Sql        sql.Config    `mapstructure:"sql"`
	Log        log.Config    `mapstructure:"log"`
	MiddleWare l.LogConfig   `mapstructure:"middleware"`
}
