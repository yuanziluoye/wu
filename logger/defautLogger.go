package logger

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/yuanziluoye/wu/config"
)

var logger = logs.NewLogger(10000)

func init() {

	appConfig := config.GetAppConfig()
	logPath := appConfig.LogPath
	logger.SetLogger("file", fmt.Sprintf(`{"filename": "%v", "perm": "0660"}`, logPath))
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(2)
	logger.SetLogger("console", fmt.Sprintf(`{"level": %d}`, logs.LevelInformational))
	logger.Info("[logger] use log path: %s", logPath)
}

func GetLogger() *logs.BeeLogger {
	return logger
}
