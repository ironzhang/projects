package command

import (
	"context"
	"time"

	"github.com/ironzhang/x-pearls/config"
	"github.com/ironzhang/zerone/zclient"
	"github.com/ironzhang/zoutil"
)

func Dial(z *zoutil.Zerone, service string) (*Client, error) {
	c, err := z.NewClient(service)
	if err != nil {
		return nil, err
	}
	return &Client{
		node:    []byte(z.Node()),
		zerone:  z,
		zclient: c.WithBalancePolicy(zclient.NodeBalancer),
	}, nil
}

type Client struct {
	node    []byte
	zerone  *zoutil.Zerone
	zclient *zclient.Client
}

func (p *Client) GetNodeName(ctx context.Context) (node string, err error) {
	err = p.zclient.Call(ctx, p.node, "command.GetNodeName", nil, &node, p.zerone.Timeout)
	return
}

func (p *Client) GetTimeout(ctx context.Context) (time.Duration, error) {
	var timeout config.Duration
	err := p.zclient.Call(ctx, p.node, "command.GetTimeout", nil, &timeout, p.zerone.Timeout)
	return time.Duration(timeout), err
}

func (p *Client) SetTimeout(ctx context.Context, timeout time.Duration) error {
	return p.zclient.Call(ctx, p.node, "command.SetTimeout", config.Duration(timeout), nil, p.zerone.Timeout)
}

func (p *Client) GetLogLevel(ctx context.Context) (level string, err error) {
	err = p.zclient.Call(ctx, p.node, "command.GetLogLevel", nil, &level, p.zerone.Timeout)
	return
}

func (p *Client) SetLogLevel(ctx context.Context, level string) error {
	return p.zclient.Call(ctx, p.node, "command.SetLogLevel", level, nil, p.zerone.Timeout)
}

func (p *Client) GetTraceVerbose(ctx context.Context, object TraceObject) (verbose int, err error) {
	err = p.zclient.Call(ctx, p.node, "command.GetTraceVerbose", object, &verbose, p.zerone.Timeout)
	return
}

func (p *Client) SetTraceVerbose(ctx context.Context, object TraceObject, verbose int) error {
	args := SetTraceVerboseArgs{TraceObject: object, Verbose: verbose}
	return p.zclient.Call(ctx, p.node, "command.SetTraceVerbose", args, nil, p.zerone.Timeout)
}
