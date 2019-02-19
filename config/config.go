package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	"gopkg.in/yaml.v2"
)

type AppWorker struct {
	Directory string   `yaml:"Directory"`
	Patterns  []string `yaml:"Patterns"`
	Command   []string `yaml:"Command"`
}

type AppConfig struct {
	Worker  []AppWorker
	Events  []string `yaml:"events"`
	LogPath string   `yaml:"logPath"`
}

var logger = logs.NewLogger(10000)
var appConfig = AppConfig{}

func init() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s \n", os.Args[0])
	}

	flag.Parse()

	loadConfig()
}

// load config
func loadConfig() {
	configFile, _ := filepath.Abs("./config.yaml")
	readData, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Error("[config] read config file failed, %v", err)
		os.Exit(0)
	}

	err = yaml.Unmarshal(readData, &appConfig)
	if err != nil {
		logger.Error("[config] parse config.yaml failed, %v", err)
		os.Exit(0)
	}

	logPath := appConfig.LogPath
	logPath, _ = filepath.Abs(logPath)
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		logger.Error("[config] create log dir failed, %v", err)
		os.Exit(0)
	}

	logger.Info("[config] load config file: %s", configFile)
}

func GetAppConfig() AppConfig {
	return appConfig
}
