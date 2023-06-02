package filewatch

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// WatchFunc 订阅回调函数
//
// 返回 true 表示停止订阅，返回 false 表示继续订阅
type WatchFunc func(path string) (quit bool)

// Watcher 文件订阅程序
type Watcher struct {
	closed   int64         // 关闭标记
	quitc    chan struct{} // 退出通道
	interval time.Duration // 观察间隔

	mu        sync.Mutex
	observers []*fileObserver
}

// NewWatcher 构建文件订阅程序
func NewWatcher(interval time.Duration) *Watcher {
	return new(Watcher).init(interval)
}

// init 初始化
func (p *Watcher) init(interval time.Duration) *Watcher {
	p.quitc = make(chan struct{})
	p.interval = interval
	go p.running()
	return p
}

// Stop 停止文件订阅程序
func (p *Watcher) Stop() {
	if atomic.CompareAndSwapInt64(&p.closed, 0, 1) {
		close(p.quitc)
	}
}

// WatchFile 订阅文件
func (p *Watcher) WatchFile(ctx context.Context, path string, f WatchFunc) {
	// 预先订阅文件
	fo := &fileObserver{path: path, watchFunc: f}
	fo.observe()

	p.mu.Lock()
	defer p.mu.Unlock()
	p.observers = append(p.observers, fo)
}

func (p *Watcher) atomicLoadObservers() []*fileObserver {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.observers
}

func (p *Watcher) executeGC() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 清理已退出的文件观察器
	observers := make([]*fileObserver, 0, len(p.observers))
	for _, fo := range p.observers {
		if fo.quit {
			continue
		}
		observers = append(observers, fo)
	}
	p.observers = observers
}

func (p *Watcher) tick() {
	observers := p.atomicLoadObservers()
	for _, fo := range observers {
		fo.observe()
	}
	p.executeGC()
}

func (p *Watcher) running() {
	t := time.NewTicker(p.interval)
	defer t.Stop()

	for {
		select {
		case <-p.quitc:
			return
		case <-t.C:
			p.tick()
		}
	}
}
