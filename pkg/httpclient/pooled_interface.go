package httpclient

import (
	middleware "github.com/Borislavv/go-httpclientpool/pkg/httpclient/middleware"
	"net/http"
)

type Pooled interface {
	Do(req *http.Request) (*http.Response, error)
	OnReq(middlewares ...middleware.RequestMiddlewareFunc) Pooled
	OnResp(middlewares ...middleware.ResponseMiddlewareFunc) Pooled
}
