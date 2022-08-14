package app

import (
	"github.com/core-go/core"
	"github.com/core-go/core/log"
	l "github.com/core-go/core/middleware/echo"
	"github.com/core-go/sql"
)

type Config struct {
	Server     core.ServerConf `mapstructure:"server"`
	Sql        sql.Config      `mapstructure:"sql"`
	Log        log.Config      `mapstructure:"log"`
	MiddleWare l.LogConfig     `mapstructure:"middleware"`
}
