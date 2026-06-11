package runner

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// EventWatcher represents a source of filesystem events.
type EventWatcher struct {
	Events chan fsnotify.Event
	Errors chan error
}

type Watcher struct {
	config  *Config
	watcher *EventWatcher
	buildFn func()
}

func NewWatcher(cfg *Config, w *EventWatcher, buildFn func()) *Watcher {
	return &Watcher{
		config:  cfg,
		watcher: w,
		buildFn: buildFn,
	}
}

func (w *Watcher) shouldIgnore(event fsnotify.Event) bool {
	base := filepath.Base(event.Name)
	if strings.HasSuffix(base, "~") {
		return true
	}
	if base == "4913" {
		return true
	}
	// hex-encoded temp files, e.g., created by some editors
	if matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8,}$`, base); matched {
		return true
	}
	return false
}

func (w *Watcher) triggerBuild() {
	if w.buildFn != nil {
		w.buildFn()
	}
}

func (w *Watcher) watch() {
	var timer *time.Timer
	delay := time.Duration(w.config.DebounceDelay) * time.Millisecond

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if !w.shouldIgnore(event) {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(delay, func() {
					w.triggerBuild()
				})
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			_ = err
		}
	}
}