package filewatch

import (
	"os"
	"time"
)

// fileObserver 文件观察器
type fileObserver struct {
	path          string    // 文件路径
	watchFunc     WatchFunc // 订阅回调函数
	lastWatchTime time.Time // 最后一次执行订阅回调函数的时间
	quit          bool      // 退出订阅标记
}

// 执行文件观察任务
func (p *fileObserver) observe() {
	fi, err := os.Stat(p.path)
	if err != nil {
		return
	}

	mt := fi.ModTime()
	if mt.After(p.lastWatchTime) {
		p.lastWatchTime = mt
		if p.watchFunc(p.path) {
			p.quit = true
		}
	}
}
