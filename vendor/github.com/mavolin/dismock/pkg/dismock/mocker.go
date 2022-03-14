package dismock

import (
	"context"
	"crypto/tls"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/session"
	"github.com/diamondburned/arikawa/state"
	"github.com/diamondburned/arikawa/utils/httputil/httpdriver"
	"github.com/stretchr/testify/assert"
)

type (
	// Mocker handles the mocking of arikawa's API calls.
	Mocker struct {
		// Server is the httptest.Server used to mock the requests.
		Server *httptest.Server
		// Client is a mocked http.Client that redirects all requests to the
		// mock-server.
		Client *http.Client
		// handlers is a map containing all handlers.
		// The outer map is sorted by path, the inner one by method.
		// This ensures that different requests don't share the same Handler array, while still
		// enforcing the call order.
		handlers map[string]map[string][]Handler // map[Path]map[HTTPMethod][]Handler
		// mut is the sync.Mutex used to secure the handlers map, when multiple request come in
		// concurrently.
		// However, mocks may not be added concurrently.
		mut *sync.Mutex
		// t is the test type called on error.
		t *testing.T
	}

	// Handler is a named handler for mocked endpoints.
	Handler struct {
		// Name is the name of the handler.
		Name string
		// Handler is the underlying http.Handler.
		http.Handler
	}

	// MockFunc is the function used to create a mock.
	MockFunc func(w http.ResponseWriter, r *http.Request, t *testing.T)
)

// New creates a new Mocker with a started server listening on
// Mocker.Server.Listener.Addr().
func New(t *testing.T) *Mocker {
	m := &Mocker{
		handlers: make(map[string]map[string][]Handler, 1),
		mut:      new(sync.Mutex),
		t:        t,
	}

	m.Server = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.mut.Lock()
		defer m.mut.Unlock()

		path := strings.TrimRight(r.URL.Path, "/")

		methHandlers, ok := m.handlers[path]
		if !assert.True(t, ok, "unhandled path '"+path+"'") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h, ok := methHandlers[r.Method]
		if !assert.True(t, ok, "unhandled method '"+r.Method+"' on path '"+path+"'") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h[0].ServeHTTP(w, r)

		if len(h) == 1 { // this is the only handler for this method
			if len(methHandlers) == 1 { // the current method is the only method for this path
				delete(m.handlers, path)
			} else { // there are other methods for this path
				delete(m.handlers[path], r.Method)
			}
		} else { // there are multiple handlers for this method
			m.handlers[path][r.Method] = m.handlers[path][r.Method][1:]
		}
	}))

	m.Server.StartTLS()

	m.Client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, _ string) (conn net.Conn, err error) {
				return net.Dial(network, m.Server.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return m
}

// NewSession creates a new Mocker, starts its test server and returns a
// manipulated session.Session using the test server.
func NewSession(t *testing.T) (*Mocker, *session.Session) {
	m := New(t)

	gw := gateway.NewCustomGateway("", "")
	s := session.NewWithGateway(gw)

	s.Client.Client.Client = httpdriver.WrapClient(*m.Client)
	s.Client.Retries = 1

	return m, s
}

// NewState creates a new Mocker, starts its test server and returns a
// manipulated state.State which's Session uses the test server.
// In order to allow for successful testing, the State's Store, will always
// return an error, forcing the use of the (mocked) Session.
func NewState(t *testing.T) (*Mocker, *state.State) {
	m, se := NewSession(t)

	st, err := state.NewFromSession(se, new(state.NoopStore))
	if err != nil {
		panic(err) // this should never happen
	}

	return m, st
}

// Mock uses the passed MockFunc to create a mock for the passed path using the
// passed method.
// If there are already handlers for this path with the same method, the
// handler will be queued up behind the other handlers with the same path and
// method.
// Queued up handlers must be invoked in the same order they were added in.
//
// Trailing slashes ('/') will be removed.
//
// Names don't need to be unique, and have the sole purpose of aiding in
// debugging.
//
// The MockFunc may be nil if only the NoContent status shall be returned.
func (m *Mocker) Mock(name, method, path string, f MockFunc) {
	path = strings.TrimRight(path, "/")

	if m.handlers[path] == nil {
		m.handlers[path] = make(map[string][]Handler)
	}

	h := Handler{
		Name: name,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if f != nil {
				f(w, r, m.t)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		}),
	}

	m.handlers[path][method] = append(m.handlers[path][method], h)
}

// MockAPI uses the passed MockFunc to as handler for the passed path and
// method.
// The path must not include the api version, i.e. '/api/v6' must be stripped.
// If there are already handlers for this path with the same method, the
// handler will be queued up behind the other handlers with the same path and
// method.
// Queued up handlers are invoked in the same order they were added in.
//
// Trailing slashes ('/') will be removed.
//
// Names don't need to be unique, and have the sole purpose of aiding in
// debugging.
//
// The MockFunc may be nil if only the NoContent status shall be returned.
func (m *Mocker) MockAPI(name, method, path string, f MockFunc) {
	path = "/api/v6" + path

	m.Mock(name, method, path, f)
}

// Clone creates a clone of the Mocker that has the same handlers but a
// separate server.
//
// Creating a clone will automatically close the Mocker's server.
func (m *Mocker) Clone(t *testing.T) (clone *Mocker) {
	m.Close()

	clone = New(t)

	clone.handlers = m.deepCopyHandlers()

	return
}

// CloneSession clones handlers of the Mocker and returns the cloned Mocker and
// a new session.Session using the new server.
//
// Creating a clone will automatically close the Mocker's server.
func (m *Mocker) CloneSession(t *testing.T) (clone *Mocker, s *session.Session) {
	m.Close()

	clone, s = NewSession(t)

	clone.handlers = m.deepCopyHandlers()

	return
}

// CloneState clones handlers of the Mocker and returns the cloned Mocker and a
// new state.State using the new server.
// Useful for multiple tests with the same API calls.
//
// Creating a clone will automatically close the current server.
func (m *Mocker) CloneState(t *testing.T) (clone *Mocker, s *state.State) {
	m.Close()

	clone, s = NewState(t)

	clone.handlers = m.deepCopyHandlers()

	return
}

// deepCopyHandlers returns a deep copy of the handlers of the Mocker.
func (m *Mocker) deepCopyHandlers() (cp map[string]map[string][]Handler) {
	cp = make(map[string]map[string][]Handler, len(m.handlers))

	for p, sub := range m.handlers {
		subCopy := make(map[string][]Handler, len(sub))

		for m, ssub := range sub {
			ssubCopy := make([]Handler, len(ssub))

			copy(ssubCopy, ssub)

			subCopy[m] = ssubCopy
		}

		cp[p] = subCopy
	}

	return
}

// Eval closes the server and evaluates if all registered handlers were
// invoked.
// If not it will call testing.T.Fatal, printing an error message with all
// uninvoked handlers.
// This must be called at the end of every test.
func (m *Mocker) Eval() {
	m.Close()

	if len(m.handlers) == 0 {
		return
	}

	m.t.Fatal("there are uninvoked handlers:\n\n" + m.genUninvokedMsg())
}

// Close shuts down the server and blocks until all current requests are
// completed.
func (m *Mocker) Close() { m.Server.Close() }

// genUninvokedMsg generates an error message stating the unused handlers.
//
// Example
//		/guilds/118456055842734083
// 			Guild: 2 uinvoked handlers
// 		/guilds/118456055842734083/members/256827968133791744
//			ModifyMember: 1 uinvoked handler
func (m *Mocker) genUninvokedMsg() string {
	missing := make(map[string]map[string]int, len(m.handlers))

	for p, methHandlers := range m.handlers {
		missing[p] = make(map[string]int, 1)

		for _, handlers := range methHandlers {
			for _, h := range handlers {
				missing[p][h.Name]++
			}
		}
	}

	n := (len(m.handlers) - 1) * 2 // 2 for the colon and the line feed

	for p, missingRequests := range missing {
		n += len(p)

		for name, qty := range missingRequests {
			// log10(qty) is the number of digits and 19 is the number of characters without
			// placeholders
			n += len(name) + int(math.Log10(float64(qty))) + 19

			if qty > 1 { // handler has to be pluralized
				n++
			}
		}
	}

	var b strings.Builder

	b.Grow(n)

	first := true

	for p, missingRequests := range missing {
		if !first {
			b.WriteRune('\n')
		}

		b.WriteString(p)
		b.WriteRune(':')

		for name, qty := range missingRequests {
			b.WriteString("\n\t")
			b.WriteString(name)
			b.WriteString(": ")
			b.WriteString(strconv.Itoa(qty))
			b.WriteString(" uninvoked handler")

			if qty > 1 {
				b.WriteRune('s')
			}
		}

		first = false
	}

	return b.String()
}
