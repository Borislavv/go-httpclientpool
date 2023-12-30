package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestPooled_Do(t *testing.T) {
	expectedResp := "fooBarBaz"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(expectedResp)); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	cfg := &Config{
		PoolInitSize: 10,
		PoolMaxSize:  1024,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, clientCancel := NewPool(ctx, cfg, func() *http.Client {
		return &http.Client{Timeout: time.Minute}
	})
	defer clientCancel()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != expectedResp {
		t.Fatalf("expected response '%v', gotten '%v'", expectedResp, string(b))
	}
}

func TestPooled_OnReq(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("")); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		PoolInitSize: 10,
		PoolMaxSize:  1024,
	}

	client, clientCancel := NewPool(context.Background(), cfg, func() *http.Client {
		return &http.Client{Timeout: time.Minute}
	})
	defer clientCancel()

	counter := &struct {
		Responses int
	}{
		Responses: 0,
	}

	client.
		OnReq(
			func(next RequestModifier) RequestModifier {
				return RequestModifierFunc(func(req *http.Request) (*http.Response, error) {
					copyValues := req.URL.Query()
					if copyValues.Has("timestamp") {
						copyValues.Del("timestamp")
					}
					copyValues.Add("timestamp", strconv.Itoa(time.Now().Nanosecond()))
					req.URL.RawQuery = copyValues.Encode()
					return next.Do(req)
				})
			},
		).
		OnResp(
			func(next ResponseHandler) ResponseHandler {
				return ResponseHandlerFunc(func(resp *http.Response, err error) (*http.Response, error) {
					counter.Responses++
					return next.Handle(resp, err)
				})
			},
			func(next ResponseHandler) ResponseHandler {
				return ResponseHandlerFunc(func(resp *http.Response, err error) (*http.Response, error) {
					counter.Responses++
					return next.Handle(resp, err)
				})
			},
		)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.Request.URL.Query().Get("timestamp") == "" {
		t.Fatal("request middleware does not applied, timestamp does not exists into the url")
	}

	if counter.Responses != 4 {
		t.Fatal("response middleware does not applied, counter is equals zero")
	}
}
