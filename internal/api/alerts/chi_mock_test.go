package alerts

import (
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type chiMock struct {
	mock.Mock
}

func (m *chiMock) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
}

func (m *chiMock) Routes() []chi.Route {
	return nil
}

func (m *chiMock) Middlewares() chi.Middlewares {
	return nil
}

func (m *chiMock) Match(rctx *chi.Context, method, path string) bool {
	args := m.Called(rctx, method, path)
	return args.Bool(0)
}

func (m *chiMock) Use(middlewares ...func(http.Handler) http.Handler) {
	m.Called(middlewares)
}
func (m *chiMock) With(middlewares ...func(http.Handler) http.Handler) chi.Router {
	args := m.Called(middlewares)
	return args.Get(0).(chi.Router)
}
func (m *chiMock) Group(fn func(r chi.Router)) chi.Router {
	args := m.Called(fn)
	return args.Get(0).(chi.Router)
}
func (m *chiMock) Route(pattern string, fn func(r chi.Router)) chi.Router {
	args := m.Called(pattern, fn)
	return args.Get(0).(chi.Router)
}
func (m *chiMock) Mount(pattern string, h http.Handler) {
	m.Called(pattern, h)
}
func (m *chiMock) Handle(pattern string, h http.Handler) {
	m.Called(pattern, h)
}
func (m *chiMock) HandleFunc(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) Method(method, pattern string, h http.Handler) {
	m.Called(method, pattern, h)
}
func (m *chiMock) MethodFunc(method, pattern string, h http.HandlerFunc) {
	m.Called(method, pattern, h)
}
func (m *chiMock) Connect(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) Delete(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) Get(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) Head(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) Options(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) Patch(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) Post(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) Put(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) Trace(pattern string, h http.HandlerFunc) {
	m.Called(pattern, h)
}
func (m *chiMock) NotFound(h http.HandlerFunc) {
	m.Called(h)
}
func (m *chiMock) MethodNotAllowed(h http.HandlerFunc) {
	m.Called(h)
}
