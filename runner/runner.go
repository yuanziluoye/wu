// Package manager provides manager for running watch and exec loop
package runner

import (
	"strings"
	"time"

	"github.com/yuanziluoye/wu/command"
	"github.com/yuanziluoye/wu/logger"
)

type Runner interface {
	Path() string
	Patterns() []string
	Command() command.Command
	Start()
	Exit()
}

type runner struct {
	path     string
	patterns []string
	command  command.Command

	abort chan struct{}
}

var appLogger = logger.GetLogger()

func New(path string, patterns []string, command command.Command) Runner {
	return &runner{
		path:     path,
		patterns: patterns,
		command:  command,
	}
}

func (r *runner) Path() string {
	return r.path
}

func (r *runner) Patterns() []string {
	return r.patterns
}

func (r *runner) Command() command.Command {
	return r.command
}

func (r *runner) Start() {
	r.abort = make(chan struct{})
	changed, err := watch(r.path, r.abort)
	if err != nil {
		appLogger.Error("Failed to initialize watcher: %v", err)
	}
	matched := match(changed, r.patterns)
	appLogger.Info("Start watching...")

	// Run the command once at initially
	r.command.Start(200 * time.Millisecond)
	for fp := range matched {
		files := gather(fp, matched, 500*time.Millisecond)

		// Terminate previous running command
		r.command.Terminate(2 * time.Second)

		appLogger.Info("File changed: %s", strings.Join(files, ", "))

		// Run new command
		r.command.Start(200 * time.Millisecond)
	}
}

func (r *runner) Exit() {
	appLogger.Info("Shutting down...")

	r.abort <- struct{}{}
	close(r.abort)
	r.command.Terminate(2 * time.Second)
}
