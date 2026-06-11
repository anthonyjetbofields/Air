package runner

import (
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestWatcherDebounce(t *testing.T) {
	cfg := &Config{
		DebounceDelay: 200,
	}

	mockEvents := make(chan fsnotify.Event)
	mockErrors := make(chan error)
	ew := &EventWatcher{
		Events: mockEvents,
		Errors: mockErrors,
	}

	var buildCount int
	var mu sync.Mutex

	buildFn := func() {
		mu.Lock()
		buildCount++
		mu.Unlock()
	}

	w := NewWatcher(cfg, ew, buildFn)

	// Run watcher in a goroutine
	go w.watch()

	// Push 5 events within a 50ms window
	for i := 0; i < 5; i++ {
		mockEvents <- fsnotify.Event{
			Name: "test.go",
			Op:   fsnotify.Write,
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for debounce period
	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	count := buildCount
	mu.Unlock()

	if count != 1 {
		t.Errorf("Expected exactly 1 build trigger, got %d", count)
	}

	// Clean up
	close(mockEvents)
	close(mockErrors)
}

func TestWatcherIgnoreSafeSave(t *testing.T) {
	cfg := &Config{
		DebounceDelay: 50,
	}

	mockEvents := make(chan fsnotify.Event)
	mockErrors := make(chan error)
	ew := &EventWatcher{
		Events: mockEvents,
		Errors: mockErrors,
	}

	var buildCount int
	var mu sync.Mutex

	buildFn := func() {
		mu.Lock()
		buildCount++
		mu.Unlock()
	}

	w := NewWatcher(cfg, ew, buildFn)

	go w.watch()

	// Push safe-save events
	mockEvents <- fsnotify.Event{Name: "file.go~", Op: fsnotify.Write}
	mockEvents <- fsnotify.Event{Name: "4913", Op: fsnotify.Write}
	mockEvents <- fsnotify.Event{Name: "a1b2c3d4", Op: fsnotify.Write} // hex encoded

	// Wait for debounce period
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	count := buildCount
	mu.Unlock()

	if count != 0 {
		t.Errorf("Expected 0 build triggers for ignored files, got %d", count)
	}

	close(mockEvents)
	close(mockErrors)
}
