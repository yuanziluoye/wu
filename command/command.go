// Package command provides a wrap over os/exec for easier command handling
package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/yuanziluoye/wu/logger"
)

type Command interface {
	String() string
	Start(delay time.Duration)
	Terminate(wait time.Duration)
}

type command struct {
	name   string
	args   []string
	cmd    *exec.Cmd
	mutex  *sync.Mutex
	exited chan struct{}
}

var appLogger = logger.GetLogger()

func New(cmdstring []string) Command {
	if len(cmdstring) == 0 {
		return Empty()
	}

	name := cmdstring[0]
	args := cmdstring[1:]

	return &command{
		name,
		args,
		nil,
		&sync.Mutex{},
		nil,
	}
}

func (c *command) String() string {
	return fmt.Sprintf("%s %s", c.name, strings.Join(c.args, " "))
}

func (c *command) Start(delay time.Duration) {
	time.Sleep(delay) // delay for a while to avoid start too frequently

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.cmd != nil && !c.cmd.ProcessState.Exited() {
		appLogger.Error("Failed to start command: previous command hasn't exit.")
	}

	cmd := exec.Command(c.name, c.args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout // Redirect stderr of sub process to stdout of parent

	// Make process group id available for the command to run
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	appLogger.Info("- Running command: %s", c.String())

	err := cmd.Start()
	exited := make(chan struct{})

	if err != nil {
		appLogger.Error("Failed: %v", err)
	} else {
		c.cmd = cmd
		c.exited = exited

		go func() {
			defer func() {
				exited <- struct{}{}
				close(exited)
			}()

			cmd.Wait()
			if cmd.ProcessState.Success() {
				appLogger.Info("- Done.")
			} else {
				appLogger.Warn("- Terminated.")
			}
		}()
	}
}

func (c *command) Terminate(wait time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// set c.cmd to nil after finished
	defer func() {
		c.cmd = nil
	}()

	if c.cmd == nil {
		// No command is runing, just return
		return
	}

	if c.cmd.ProcessState != nil && c.cmd.ProcessState.Exited() {
		// Command has exited, just return
		return
	}

	appLogger.Info("- Stopping")
	// Try to stop the process by sending a SIGINT signal
	if err := c.kill(syscall.SIGINT); err != nil {
		appLogger.Error("Failed to terminate process with interrupt: %v", err)
	}

	for {
		select {
		case <-c.exited:
			return
		case <-time.After(wait):
			appLogger.Warn("- Killing process")
			c.kill(syscall.SIGTERM)
		}
	}
}

func (c *command) kill(sig syscall.Signal) error {
	cmd := c.cmd
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		return syscall.Kill(-pgid, sig)
	}
	return err
}
