package hapi

import (
    "net/http"

    "github.com/gorilla/context"
)

// A common function for various test cases
func (h *HypermediaAPI) TestRegister(ctype,id string) {
    h.Register("GET", "/", ctype, func(w http.ResponseWriter, r *http.Request) { context.Set(r,"id",id) })
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
    return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
    return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
    return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}
