package file

import "time"

type timer struct {
	running bool
	*time.Timer
}

func (t *timer) Init(d time.Duration, running bool) {
	t.running = running
	t.Timer = time.NewTimer(d)
	if !running {
		t.Timer.Stop()
	}
}

func (t *timer) Stop() bool {
	t.running = false
	return t.Timer.Stop()
}

func (t *timer) Reset(d time.Duration) bool {
	t.running = true
	return t.Timer.Reset(d)
}

func (t *timer) Running() bool {
	return t.running
}
