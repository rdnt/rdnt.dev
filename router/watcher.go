package router

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

type Watcher struct {
	watcher     *watcher.Watcher
	interval    time.Duration
	handlerFunc func(string, []byte)
	journalRex  *regexp.Regexp
}

func NewWatcher(h func(string, []byte), d time.Duration) (*Watcher, error) {
	return &Watcher{
		watcher:     watcher.New(),
		interval:    d,
		handlerFunc: h,
	}, nil
}

func (w *Watcher) Watch(path string) error {
	err := w.watcher.AddRecursive(path)
	if err != nil {
		return err
	}
	return nil
}

func (w *Watcher) Read(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}

	w.handlerFunc(path, b)
}

func (w *Watcher) Start() {
	fmt.Println("starting watcher")
	go func() {
		for {
			select {
			case event := <-w.watcher.Event:
				go w.Read(event.Path) // Print the event's info.
			case err := <-w.watcher.Error:
				fmt.Println(err)
				continue
			case <-w.watcher.Closed:
				return
			}
		}
	}()

	go func() {
		_ = w.watcher.Start(w.interval)
	}()
}

func (w *Watcher) Stop() {
	w.watcher.Close()
}
