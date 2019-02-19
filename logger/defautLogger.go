package logger

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/yuanziluoye/wu/config"
)

var logger = logs.NewLogger(10000)

func init() {

	appConfig := config.GetAppConfig()
	logPath := appConfig.Logger.LogPath
	loggerConfig := appConfig.Logger
	loggerConfigJson := `{"filename":"%v", "daily": %v, "maxDays": %v, "rotate": %v, "level": %v, "perm":"%v", "rotateperm":"%v"}`
	loggerJsonConfig := fmt.Sprintf(loggerConfigJson, logPath, loggerConfig.Daily, loggerConfig.MaxDays,
		loggerConfig.Rotate, loggerConfig.Level, loggerConfig.Perm, loggerConfig.RotatePerm)

	logger.SetLogger("file", loggerJsonConfig)
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(2)
	logger.SetLogger("console", fmt.Sprintf(`{"level": %d}`, logs.LevelInformational))

	logger.Info("[init] use log path: %s", logPath)
}

func GetLogger() *logs.BeeLogger {
	return logger
}
