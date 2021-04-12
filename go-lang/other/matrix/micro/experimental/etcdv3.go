package experimental

import (
	"context"

	"github.com/coreos/etcd/clientv3"
)

func register(ctx context.Context, c *clientv3.Client, key, value string, ttl int64) (clientv3.LeaseID, error) {
	resp, err := c.Grant(ctx, ttl)
	if err != nil {
		return clientv3.NoLease, err
	}
	if _, err = c.Put(ctx, key, value, clientv3.WithLease(resp.ID)); err != nil {
		return clientv3.NoLease, err
	}
	return resp.ID, nil
}

func unregister(ctx context.Context, c *clientv3.Client, id clientv3.LeaseID, key string) (err error) {
	c.Delete(ctx, key)
	if _, err = c.Revoke(ctx, id); err != nil {
		return err
	}
	return nil
}
