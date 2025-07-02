package gatekeeping

import (
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

func NewMatcher() *Matcher {
	return &Matcher{router: httprouter.New()}
}

func (m *Matcher) Load(endpoints []Endpoint) {
	m.lock.Lock()
	defer m.lock.Unlock()

	r := httprouter.New()
	for _, ep := range endpoints {
		code := ep.Code
		method := ep.Method
		r.Handle(method, ep.Path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.Header().Set("x-endpoint-code", code)
		})
	}

	m.router = r
}

func (m *Matcher) Add(e Endpoint) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.router.Handle(e.Method, e.Path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("x-endpoint-code", e.Code)
	})
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
func (c *capture) handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c.found = true
	c.code = w.Header().Get("x-endpoint-code")
}
