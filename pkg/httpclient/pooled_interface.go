package httpclient

import (
	"net/http"
)

type Pooled interface {
	Do(req *http.Request) (*http.Response, error)
	OnReq(middlewares ...RequestMiddlewareFunc) Pooled
	OnResp(middlewares ...ResponseMiddlewareFunc) Pooled
}
