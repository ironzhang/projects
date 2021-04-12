package restful

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/ironzhang/matrix/codes"
	"github.com/ironzhang/matrix/context-value"
	"github.com/ironzhang/matrix/httputils"
	"github.com/ironzhang/matrix/restful/codec"
	"github.com/ironzhang/matrix/tlog"
	"github.com/ironzhang/matrix/uuid"
)

var DefaultClient = &Client{
	Client: &http.Client{
		Transport: httputils.NewVerboseRoundTripper(nil, nil, nil),
		Timeout:   20 * time.Second,
	},
}

type Client struct {
	Client  *http.Client
	Codec   codec.Codec
	Context context.Context
}

func (c *Client) Delete(url string, args, reply interface{}) error {
	return c.Do("DELETE", url, args, reply)
}

func (c *Client) Get(url string, args, reply interface{}) error {
	return c.Do("GET", url, args, reply)
}

func (c *Client) Head(ctx context.Context, url string, args, reply interface{}) error {
	return c.Do("HEAD", url, args, reply)
}

func (c *Client) Options(ctx context.Context, url string, args, reply interface{}) error {
	return c.Do("OPTIONS", url, args, reply)
}

func (c *Client) Patch(url string, args, reply interface{}) error {
	return c.Do("PATCH", url, args, reply)
}

func (c *Client) Post(url string, args, reply interface{}) error {
	return c.Do("POST", url, args, reply)
}

func (c *Client) Put(url string, args, reply interface{}) error {
	return c.Do("PUT", url, args, reply)
}

func (c *Client) Do(method, url string, args, reply interface{}) error {
	return c.DoContext(c.context(), method, url, args, reply)
}

func (c *Client) DeleteContext(ctx context.Context, url string, args, reply interface{}) error {
	return c.DoContext(ctx, "DELETE", url, args, reply)
}

func (c *Client) GetContext(ctx context.Context, url string, args, reply interface{}) error {
	return c.DoContext(ctx, "GET", url, args, reply)
}

func (c *Client) HeadContext(ctx context.Context, url string, args, reply interface{}) error {
	return c.DoContext(ctx, "HEAD", url, args, reply)
}

func (c *Client) OptionsContext(ctx context.Context, url string, args, reply interface{}) error {
	return c.DoContext(ctx, "OPTIONS", url, args, reply)
}

func (c *Client) PatchContext(ctx context.Context, url string, args, reply interface{}) error {
	return c.DoContext(ctx, "PATCH", url, args, reply)
}

func (c *Client) PostContext(ctx context.Context, url string, args, reply interface{}) error {
	return c.DoContext(ctx, "POST", url, args, reply)
}

func (c *Client) PutContext(ctx context.Context, url string, args, reply interface{}) error {
	return c.DoContext(ctx, "PUT", url, args, reply)
}

func (c *Client) DoContext(ctx context.Context, method, url string, args, reply interface{}) (err error) {
	ctx = contextWithTraceId(ctx)
	log := tlog.WithContext(ctx).Sugar().With("method", method, "url", url)

	var b bytes.Buffer

	// Encode
	if args != nil {
		if err = c.codec().Encode(&b, args); err != nil {
			log.Errorw("encode", "error", err)
			return err
		}
	}

	// New http request
	req, err := http.NewRequest(method, url, &b)
	if err != nil {
		log.Errorw("new request", "error", err)
		return err
	}
	c.setHeader(ctx, req.Header)

	// Do
	resp, err := c.client().Do(req)
	if err != nil {
		log.Errorw("client do", "error", err)
		return err
	}
	defer resp.Body.Close()

	// Handle error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var e codec.Error
		if err = c.codec().DecodeError(resp.Body, &e); err != nil {
			log.Errorw("decode error", "error", err, "status", http.StatusText(resp.StatusCode))
			return err
		}
		return Errorf(resp.StatusCode, codes.Code(e.Code), e.Cause)
	}

	// Decode
	if reply != nil {
		if err = c.codec().Decode(resp.Body, reply); err != nil {
			log.Errorw("decode", "error", err, "status", http.StatusText(resp.StatusCode))
			return err
		}
	}

	return nil
}

func (c *Client) setHeader(ctx context.Context, h http.Header) {
	h.Set("Content-Type", c.codec().ContentType())
	if v := context_value.ParseTraceId(ctx); v != "" {
		h.Set(httputils.X_TRACE_ID, v)
	}
	if v := context_value.ParseVerbose(ctx); v {
		h.Set(httputils.X_VERBOSE, "1")
	}
}

func (c *Client) client() *http.Client {
	if c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}

func (c *Client) codec() codec.Codec {
	if c.Codec == nil {
		return codec.DefaultCodec
	}
	return c.Codec
}

func (c *Client) context() context.Context {
	if c.Context == nil {
		return context.Background()
	}
	return c.Context
}

func contextWithTraceId(ctx context.Context) context.Context {
	if v := context_value.ParseTraceId(ctx); v == "" {
		return context_value.WithTraceId(ctx, uuid.New().String())
	}
	return ctx
}
