package gatekeeping

import (
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

func NewMatcher() *Matcher {
	return &Matcher{
		router:    httprouter.New(),
		endpoints: make(map[string]Endpoint),
	}
}

func (m *Matcher) Load(endpoints []Endpoint) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.endpoints = make(map[string]Endpoint)
	r := httprouter.New()

	for _, ep := range endpoints {
		m.endpoints[ep.Code] = ep

		ep := ep // capture variable for closure
		r.Handle(ep.Method, ep.Path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.Header().Set("x-endpoint-code", ep.Code)
		})
	}

	m.router = r
}

func (m *Matcher) Add(e Endpoint) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.endpoints[e.Code] = e

	ep := e // capture variable for closure
	m.router.Handle(ep.Method, ep.Path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("x-endpoint-code", ep.Code)
	})
}

func (m *Matcher) Drop(code string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.endpoints, code)

	r := httprouter.New()
	for _, ep := range m.endpoints {
		ep := ep // capture variable for closure
		r.Handle(ep.Method, ep.Path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.Header().Set("x-endpoint-code", ep.Code)
		})
	}
	m.router = r
}

func (m *Matcher) Match(method, path string) (string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	rw := &capture{header: http.Header{}}
	req := &http.Request{
		Method: method,
		URL:    &url.URL{Path: path},
	}

	m.router.ServeHTTP(rw, req)

	return rw.code, rw.found
}

func (c *capture) Header() http.Header         { return c.header }
func (c *capture) Write(_ []byte) (int, error) { return 0, nil }
func (c *capture) WriteHeader(_ int)           {}

// Optional helper to extract header after Match
func (c *capture) handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c.found = true
	c.code = w.Header().Get("x-endpoint-code")
}
