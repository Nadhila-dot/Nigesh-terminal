package terminal

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Loader struct {
	frames  []string
	message string
	active  bool
	stopCh  chan struct{}
	doneCh  chan struct{}
	mu      sync.Mutex
}

func NewLoader() *Loader {
	return &Loader{
		frames: []string{"■░░░", "█■░░", "░█■░", "░░█■", "░░░█", "■░░█", "█■░█", "░█■░", "░░█■"},
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
}

func (l *Loader) Set(msg string) {
	l.mu.Lock()
	l.message = msg
	l.mu.Unlock()
	if !l.active {
		l.active = true
		go l.start()
	}
}

func (l *Loader) start() {
	i := 0
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		fmt.Print("\r\033[K") // clear line
		close(l.doneCh)
	}()
	for {
		select {
		case <-l.stopCh:
			return
		case <-sigCh:
			return
		default:
			l.mu.Lock()
			msg := l.message
			l.mu.Unlock()
			fmt.Printf("\r \033[K%s %s", l.frames[i%len(l.frames)], msg) // clear line before printing
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}

func (l *Loader) Stop() {
	if l.active {
		close(l.stopCh)
		<-l.doneCh
		l.active = false
	}
	fmt.Print("\r\033[K") // clear line
}

func (l *Loader) Exit() {
	l.Stop()
}
