package initialize

import (
	"path/filepath"

	"github.com/issueye/version-mana/internal/config"
	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/pkg/logger"
)

func InitLog() {
	logConf := new(logger.Config)
	logConf.Path = filepath.Join("runtime", "logs")
	logConf.MaxSize = config.GetParam("LOG-MAX-SIZE", "10").Int()
	logConf.MaxBackups = config.GetParam("LOG-MAX-BACKUPS", "10").Int()
	logConf.MaxAge = config.GetParam("LOG-MAX-AGE", "10").Int()
	logConf.Compress = config.GetParam("LOG-COMPRESS", "10").Bool()
	logConf.Level = config.GetParam("LOG-LEVEL", "10").Int()
	global.Log, global.Logger = logger.InitLogger(logConf)
}
