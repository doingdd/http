// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package http

import (
	"context"
	nethttp "net/http"
	"net/http/httptest"
	"strings"

	api "github.com/gdt-dev/gdt/api"
)

const (
	StateKeyBaseURL = "http.base_url"
	StateKeyClient  = "http.client"
)

type httpServerFixture struct {
	handler nethttp.Handler
	server  *httptest.Server
	useTLS  bool
}

func (f *httpServerFixture) Start(ctx context.Context) error {
	if !f.useTLS {
		f.server = httptest.NewServer(f.handler)
	} else {
		f.server = httptest.NewTLSServer(f.handler)
	}
	return nil
}

func (f *httpServerFixture) Stop(ctx context.Context) {
	f.server.Close()
}

func (f *httpServerFixture) HasState(key string) bool {
	lkey := strings.ToLower(key)
	switch lkey {
	case StateKeyBaseURL, StateKeyClient:
		return true
	}
	return false
}

func (f *httpServerFixture) State(key string) interface{} {
	key = strings.ToLower(key)
	switch key {
	case StateKeyBaseURL:
		return f.server.URL
	case StateKeyClient:
		return f.server.Client()
	}
	return ""
}

// NewServerFixture returns a fixture that will start and stop a supplied
// http.Handler. The returned fixture exposes an "http.base_url" state key that
// test cases of type "http" examine to determine the base URL the tests should
// hit
func NewServerFixture(h nethttp.Handler, useTLS bool) api.Fixture {
	return &httpServerFixture{handler: h, useTLS: useTLS}
}
