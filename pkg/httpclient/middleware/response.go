package httpclientmiddleware

import "net/http"

type ResponseMiddleware interface {
	Exec(next ResponseHandler) ResponseHandler
}

type ResponseMiddlewareFunc func(next ResponseHandler) ResponseHandler

func (m ResponseMiddlewareFunc) Exec(next ResponseHandler) ResponseHandler {
	return m(next)
}

type ResponseHandler interface {
	Handle(resp *http.Response, err error) (*http.Response, error)
}

type ResponseHandlerFunc func(resp *http.Response, err error) (*http.Response, error)

func (m ResponseHandlerFunc) Handle(resp *http.Response, err error) (*http.Response, error) {
	return m(resp, err)
}
