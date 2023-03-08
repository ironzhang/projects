package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/tlog"
)

func newPinger(c *clientv3.Client, timeout time.Duration, ttl int64, key, value string) *pinger {
	return &pinger{
		client:  c,
		timeout: timeout,
		ttl:     ttl,
		key:     key,
		value:   value,
		done:    make(chan struct{}),
		ok:      make(chan struct{}),
	}
}

type pinger struct {
	client  *clientv3.Client
	timeout time.Duration
	ttl     int64
	key     string
	value   string

	done chan struct{}
	ok   chan struct{}

	leaseID clientv3.LeaseID
}

func (p *pinger) Setup() error {
	log := tlog.Std().Sugar().With("key", p.key, "value", p.value)

	// key exist?
	exist, err := p.Exist(context.Background())
	if err != nil {
		log.Errorw("exist", "error", err)
		return err
	}
	if exist {
		log.Errorw("key existed")
		return fmt.Errorf("key(%s) existed", p.key)
	}

	// grant
	if err = p.Grant(context.Background()); err != nil {
		log.Errorw("grant", "error", err)
		return err
	}

	// put
	if err = p.Put(context.Background()); err != nil {
		log.Errorw("put", "error", err)
		return err
	}

	// go
	p.Go()

	log.Debug("pinger setup")
	return nil
}

func (p *pinger) Close() (err error) {
	log := tlog.Std().Sugar().With("key", p.key, "value", p.value)

	// quit run
	close(p.done)
	<-p.ok

	// delete
	if err = p.Delete(context.Background()); err != nil {
		log.Errorw("delete", "error", err)
		return err
	}

	// revoke
	if err = p.Revoke(context.Background()); err != nil {
		log.Errorw("revoke", "error", err)
		return err
	}

	log.Debug("pinger close")
	return nil
}

func (p *pinger) Go() {
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-p.done
			cancel()
		}()

		p.Run(ctx)
		close(p.ok)
	}()
}

func (p *pinger) Run(ctx context.Context) {
	log := tlog.Std().Sugar().With("key", p.key, "value", p.value)

	t := time.NewTicker(time.Duration(p.ttl) * time.Second)
	defer t.Stop()

	var err error
	var exist bool
	for {
		select {
		case <-t.C:
			if exist, err = p.Exist(ctx); err != nil {
				log.Errorw("exist", "error", err)
				continue
			}
			if exist {
				continue
			}
			if err = p.Grant(ctx); err != nil {
				log.Errorw("grant", "error", err)
				continue
			}
			if err = p.Put(ctx); err != nil {
				log.Errorw("put", "error", err)
				continue
			}
			log.Debug("ping")
		case <-ctx.Done():
			log.Debug("pinger quit")
			return
		}
	}
}

func (p *pinger) Grant(ctx context.Context) error {
	// grant
	resp, err := p.client.Grant(p.WithTimeout(ctx), p.ttl)
	if err != nil {
		return err
	}

	// keep alive, note: don't with timeout
	if _, err = p.client.KeepAlive(context.Background(), resp.ID); err != nil {
		return err
	}

	p.leaseID = resp.ID
	return nil
}

func (p *pinger) Revoke(ctx context.Context) error {
	if _, err := p.client.Revoke(p.WithTimeout(ctx), p.leaseID); err != nil {
		return err
	}
	return nil
}

func (p *pinger) Put(ctx context.Context) error {
	if _, err := p.client.Put(p.WithTimeout(ctx), p.key, p.value, clientv3.WithLease(p.leaseID)); err != nil {
		return err
	}
	return nil
}

func (p *pinger) Delete(ctx context.Context) error {
	if _, err := p.client.Delete(p.WithTimeout(ctx), p.key); err != nil {
		return err
	}
	return nil
}

func (p *pinger) Exist(ctx context.Context) (bool, error) {
	resp, err := p.client.Get(p.WithTimeout(ctx), p.key)
	if err != nil {
		return false, err
	}
	if len(resp.Kvs) == 0 {
		return false, nil
	}
	return true, nil
}

func (p *pinger) WithTimeout(ctx context.Context) context.Context {
	if p.timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, p.timeout)
	}
	return ctx
}
