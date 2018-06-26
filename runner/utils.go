package runner

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yuanziluoye/wu/config"
	"github.com/yuanziluoye/wu/logger"
)

var appConfig = config.GetAppConfig()
var appLoggerSecond = logger.GetLogger()

func watch(path string, abort <-chan struct{}) (<-chan string, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for p := range list(path) {
		err = watcher.Add(p)
		if err != nil {
			appLoggerSecond.Error("Failed to watch: %s, error: %s", p, err)
		}
	}

	out := make(chan string)
	go func() {
		defer close(out)
		defer watcher.Close()
		for {
			select {
			case <-abort:
				// Abort watching
				err := watcher.Close()
				if err != nil {
					appLoggerSecond.Error("Failed to stop watch")
				}
				return
			case fp := <-watcher.Events:

				if !itemInSlice(fp.Op.String(), appConfig.Events) {
					break
				}

				if fp.Op == fsnotify.Create {
					info, err := os.Stat(fp.Name)
					if err == nil && info.IsDir() {
						// Add newly created sub directories to watch list
						watcher.Add(fp.Name)
					}
				}
				out <- fp.Name
			case err := <-watcher.Errors:
				appLoggerSecond.Error("Watch Error: %v", err)
			}
		}
	}()

	return out, nil
}

func match(in <-chan string, patterns []string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for fp := range in {
			info, err := os.Stat(fp)
			if os.IsNotExist(err) || !info.IsDir() {
				_, fn := filepath.Split(fp)
				for _, p := range patterns {
					if ok, _ := filepath.Match(p, fn); ok {
						out <- fp
					}
				}
			}
		}
	}()

	return out
}

func list(root string) <-chan string {
	out := make(chan string)

	info, err := os.Stat(root)
	if err != nil {
		appLoggerSecond.Error("Failed to visit %s, error: %s\n", root, err)
	}
	if !info.IsDir() {
		go func() {
			defer close(out)
			out <- root
		}()

		return out
	}

	go func() {
		defer close(out)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if err != nil {
					appLoggerSecond.Error("Failed to visit directory: %s, error: %s", path, err)
					return err
				}
				out <- path
			}
			return nil
		})
	}()

	return out
}

// gather delays further operations for a while and gather
// all changes happened in this period
func gather(first string, changes <-chan string, delay time.Duration) []string {
	files := make(map[string]bool)
	files[first] = true
	after := time.After(delay)
loop:
	for {
		select {
		case fp := <-changes:
			files[fp] = true
		case <-after:
			// After the delay, return collected filenames
			break loop
		}
	}

	ret := []string{}
	for k := range files {
		ret = append(ret, k)
	}

	sort.Strings(ret)
	return ret
}

func itemInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
