package app

import (
	"github.com/core-go/core"
	"github.com/core-go/log"
	"github.com/core-go/sql"

	l "go-service/pkg/log"
)

type Config struct {
	Server     core.ServerConf `mapstructure:"server"`
	Sql        sql.Config      `mapstructure:"sql"`
	Log        log.Config      `mapstructure:"log"`
	MiddleWare l.LogConfig     `mapstructure:"middleware"`
}
