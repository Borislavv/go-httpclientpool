package httpclient

import (
	"context"
	config "ddos/internal/ddos/infrastructure/httpclient/config"
	middleware "ddos/internal/ddos/infrastructure/httpclient/middleware"
	"net/http"
)

type CancelFunc func()

type Pool struct {
	ctx     context.Context
	pool    chan *http.Client
	creator func() *http.Client
	cancel  context.CancelFunc
	req     middleware.RequestModifier
	resp    middleware.ResponseHandler
}

func NewPool(ctx context.Context, cfg *config.Config, creator func() *http.Client) (*Pool, CancelFunc) {
	if cfg.PoolInitSize > cfg.PoolMaxSize {
		cfg.PoolInitSize = cfg.PoolMaxSize
	}

	ctx, cancel := context.WithCancel(ctx)

	p := &Pool{
		ctx:     ctx,
		pool:    make(chan *http.Client, cfg.PoolMaxSize),
		creator: creator,
		cancel:  cancel,
	}

	for i := 0; i < cfg.PoolInitSize; i++ {
		p.pool <- creator()
	}

	p.req = middleware.RequestModifierFunc(
		func(req *http.Request) (*http.Response, error) {
			c := p.get()
			defer p.put(c)
			return c.Do(req)
		},
	)

	p.resp = middleware.ResponseHandlerFunc(
		func(resp *http.Response, err error) (*http.Response, error) {
			return resp, err
		},
	)

	return p, p.cls
}

func (p *Pool) Do(req *http.Request) (*http.Response, error) {
	return p.resp.Handle(p.req.Do(req.WithContext(p.ctx)))
}

func (p *Pool) OnReq(middlewares ...middleware.RequestMiddlewareFunc) Pooled {
	for i := len(middlewares) - 1; i >= 0; i-- {
		p.req = middlewares[i].Exec(p.req)
	}
	return p
}

func (p *Pool) OnResp(middlewares ...middleware.ResponseMiddlewareFunc) Pooled {
	for i := len(middlewares) - 1; i >= 0; i-- {
		p.resp = middlewares[i].Exec(p.resp)
	}
	return p
}

func (p *Pool) get() *http.Client {
	select {
	case c := <-p.pool:
		return c
	default:
		return p.creator()
	}
}

func (p *Pool) put(c *http.Client) {
	select {
	case p.pool <- c:
	default:
	}
}

func (p *Pool) cls() {
	p.cancel()
	close(p.pool)
	p.pool = nil
}
