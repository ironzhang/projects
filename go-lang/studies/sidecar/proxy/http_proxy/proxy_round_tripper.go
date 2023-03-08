package http_proxy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ironzhang/tlog"
	"go.uber.org/multierr"
)

// AddressSelector 地址选择器
type AddressSelector interface {
	SelectAddress(ctx context.Context, tags map[string]string) (string, error)
}

// AddressStateMarker 地址状态标记器
type AddressStateMarker interface {
	MarkAddressState(ctx context.Context, addr string, fault bool)
}

type proxyRoundTripper struct {
	selector  AddressSelector
	transport http.Transport
}

func (p *proxyRoundTripper) markAddressState(ctx context.Context, addr string, err error) {
	if m, ok := p.selector.(AddressStateMarker); ok {
		if err != nil {
			m.MarkAddressState(ctx, addr, true)
		} else {
			m.MarkAddressState(ctx, addr, false)
		}
	}
}

func (p *proxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	tags := make(map[string]string)
	tries := 2

	var merr error
	for i := 0; i < tries; i++ {
		addr, err := p.selector.SelectAddress(ctx, tags)
		if err != nil {
			merr = multierr.Append(merr, err)
			tlog.WithContext(ctx).Warnw("select address", "error", err)
			break
		}

		req.URL.Host = addr
		resp, err := p.transport.RoundTrip(req)
		p.markAddressState(ctx, addr, err)
		if err != nil {
			merr = multierr.Append(merr, err)
			tlog.WithContext(ctx).Warnw("round trip", "url_host", req.URL.Host, "url_path", req.URL.Path, "error", err)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("proxy round trip: %w", merr)
}

func (p *proxyRoundTripper) parseRequestTags(req *http.Request) map[string]string {
	return make(map[string]string)
}
