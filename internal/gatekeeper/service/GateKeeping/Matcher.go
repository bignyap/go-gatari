package gatekeeping

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// === Constructor ===

func NewMatcher() *Matcher {
	return &Matcher{
		router:    httprouter.New(),
		endpoints: make(map[string]Endpoint),
	}
}

// === Helpers ===

func normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func routeKey(method, path string) string {
	return strings.ToUpper(method) + " " + normalizePath(path)
}

// === Load (replace all) ===

func (m *Matcher) Load(endpoints []Endpoint) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.endpoints = make(map[string]Endpoint)
	r := httprouter.New()

	for _, ep := range endpoints {
		ep.Path = normalizePath(ep.Path)
		key := routeKey(ep.Method, ep.Path)

		if _, exists := m.endpoints[key]; exists {
			fmt.Printf("⚠️ Skipping duplicate: [%s] %s (code: %s)\n", ep.Method, ep.Path, ep.Code)
			continue
		}

		m.endpoints[key] = ep

		ep := ep // capture for closure
		r.Handle(ep.Method, ep.Path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.Header().Set("x-endpoint-code", ep.Code)
			w.WriteHeader(http.StatusOK)
		})
		fmt.Printf("✔️ Registered: [%s] %s => Code: %s\n", ep.Method, ep.Path, ep.Code)
	}

	m.router = r
}

// === Add (idempotent) ===

func (m *Matcher) Add(e Endpoint) {
	m.lock.Lock()
	defer m.lock.Unlock()

	e.Path = normalizePath(e.Path)
	key := routeKey(e.Method, e.Path)

	if _, exists := m.endpoints[key]; exists {
		return // Already registered
	}

	m.endpoints[key] = e

	ep := e
	m.router.Handle(ep.Method, ep.Path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("x-endpoint-code", ep.Code)
		w.WriteHeader(http.StatusOK)
	})
}

// === Drop (by code, rebuilds router) ===

func (m *Matcher) Drop(code string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for key, ep := range m.endpoints {
		if ep.Code == code {
			delete(m.endpoints, key)
		}
	}

	r := httprouter.New()
	for _, ep := range m.endpoints {
		ep := ep
		r.Handle(ep.Method, ep.Path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.Header().Set("x-endpoint-code", ep.Code)
			w.WriteHeader(http.StatusOK)
		})
	}
	m.router = r
}

// === Match ===

func (m *Matcher) Match(method, path string) (string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	rw := &capture{header: http.Header{}}
	req := &http.Request{
		Method: method,
		URL:    &url.URL{Path: normalizePath(path)},
	}

	m.router.ServeHTTP(rw, req)
	return rw.code, rw.found
}

// === capture implementation ===

func (c *capture) Header() http.Header         { return c.header }
func (c *capture) Write(_ []byte) (int, error) { return 0, nil }

func (c *capture) WriteHeader(_ int) {
	c.found = true
	c.code = c.header.Get("x-endpoint-code")
}
