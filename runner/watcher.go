func (w *Watcher) watch() {
	var timer *time.Timer
	delay := time.Duration(w.config.DebounceDelay) * time.Millisecond

	for {
		select {
		case event := <-w.watcher.Events:
			if !w.shouldIgnore(event) {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(delay, func() {
					w.triggerBuild()
				})
			}
		case err := <-w.watcher.Errors:
			// handle errors
		}
	}
}