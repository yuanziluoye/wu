// Package main provides entry for the command line tool
package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"time"
	"syscall"

	"github.com/yuanziluoye/wu/command"
	"github.com/yuanziluoye/wu/runner"
	"github.com/yuanziluoye/wu/config"
	"github.com/yuanziluoye/wu/logger"
)

var appConfig = config.GetAppConfig()

var appLogger = logger.GetLogger()

func main() {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)

	appLogger.Notice("[main] app starting, %s", time.Now().Format("2006-01-02 15:04:05.999999999"))

	quitChan := make(chan bool)
	workerLength := len(appConfig.Worker)
	workRunners := make([]runner.Runner, workerLength)

	for i := 0; i < workerLength; i++ {

		workerDirectory := parseDirectory(appConfig.Worker[i].Directory)
		patterns := appConfig.Worker[i].Patterns
		workerCommand := appConfig.Worker[i].Command

		absPath, _ := filepath.Abs(workerDirectory)
		cmd := command.New(workerCommand)

		workRunners[i] = runner.New(absPath, patterns, cmd)

		go func(i int) {
			workRunners[i].Start()
		}(i)
	}

	go func() {
		sig := <-signals
		appLogger.Warning("[signal] received sig %s, quit", sig.String())

		for _, oneRunner := range workRunners {
			oneRunner.Exit()
		}

		quitChan <- true

	}()

	<-quitChan

	appLogger.Notice("[main] app exited, %s", time.Now().Format("2006-01-02 15:04:05.999999999"))
	appLogger.Flush()
}

func parseDirectory(dir string) string {

	if info, err := os.Stat(dir); err == nil {
		if !info.IsDir() {
			appLogger.Error("[main] %v is not a directory", dir)
			os.Exit(0)
		}
	}
	return dir
}
