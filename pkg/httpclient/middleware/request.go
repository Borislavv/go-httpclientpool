package httpclientmiddleware

import "net/http"

type RequestMiddleware interface {
	Exec(next RequestModifier) RequestModifier
}

type RequestMiddlewareFunc func(next RequestModifier) RequestModifier

func (m RequestMiddlewareFunc) Exec(next RequestModifier) RequestModifier {
	return m(next)
}

type RequestModifier interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestModifierFunc func(req *http.Request) (*http.Response, error)

func (m RequestModifierFunc) Do(req *http.Request) (*http.Response, error) {
	return m(req)
}
