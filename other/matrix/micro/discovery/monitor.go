package discovery

import (
	"context"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ironzhang/matrix/tlog"
)

type monitor struct {
	client   *clientv3.Client
	revision int64
	prefix   string
	kvs      map[string][]byte
	refresh  func(kvs map[string][]byte)
}

func (m *monitor) Setup(timeout time.Duration) error {
	log := tlog.Std().Sugar().With("prefix", m.prefix)

	ctx := context.Background()
	if timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, timeout)
	}

	resp, err := m.client.Get(ctx, m.prefix, clientv3.WithPrefix())
	if err != nil {
		log.Errorw("get", "error", err)
		return err
	}
	m.revision = resp.Header.Revision
	for _, kv := range resp.Kvs {
		m.Put(kv)
	}
	log.Debugw("monitor setup", "kvs", m.kvs)

	m.Refresh()
	return nil
}

func (m *monitor) Go(done <-chan struct{}) <-chan struct{} {
	ok := make(chan struct{})
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-done
			cancel()
		}()

		m.Run(ctx)
		close(ok)
	}()
	return ok
}

func (m *monitor) Run(ctx context.Context) {
	log := tlog.Std().Sugar().With("prefix", m.prefix)
	watchc := m.Watch(ctx)
	for {
		select {
		case resp, ok := <-watchc:
			if !ok {
				watchc = m.Watch(ctx)
				continue
			}
			m.revision = resp.Header.Revision
			for _, event := range resp.Events {
				m.Update(event)
			}
			m.Refresh()
		case <-ctx.Done():
			log.Debug("monitor quit")
			return
		}
	}
}

func (m *monitor) Watch(ctx context.Context) clientv3.WatchChan {
	var revision int64
	if m.revision != 0 {
		revision = m.revision + 1
	}
	tlog.Std().Sugar().Debugw("monitor watch", "prefix", m.prefix, "revision", revision)
	return m.client.Watch(ctx, m.prefix, clientv3.WithPrefix(), clientv3.WithRev(revision))
}

func (m *monitor) Update(e *clientv3.Event) {
	switch e.Type {
	case clientv3.EventTypePut:
		m.Put(e.Kv)
	case clientv3.EventTypeDelete:
		m.Delete(e.Kv)
	}
}

func (m *monitor) Put(kv *mvccpb.KeyValue) {
	key := strings.TrimPrefix(string(kv.Key), m.prefix)
	m.kvs[key] = kv.Value
	tlog.Std().Sugar().Debugw("monitor put", "prefix", m.prefix, "key", key, "value", string(kv.Value))
}

func (m *monitor) Delete(kv *mvccpb.KeyValue) {
	key := strings.TrimPrefix(string(kv.Key), m.prefix)
	delete(m.kvs, key)
	tlog.Std().Sugar().Debugw("monitor delete", "prefix", m.prefix, "key", key)
}

func (m *monitor) Refresh() {
	if m.refresh != nil {
		m.refresh(m.kvs)
	}
}
