package hapi

import (
    "testing"
    "net/http"
)

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

func TestRouting(t *testing.T) {
    router := New()
    handlerID := ""
    router.Register("GET","/","text/html", func (w http.ResponseWriter, r *http.Request, p Params) { handlerID = "1" })
    router.UnsupportedMediaType = func (w http.ResponseWriter, r *http.Request, p Params) { handlerID = "unsupported" }
    w := new(mockResponseWriter)
    req,_ := http.NewRequest("GET","/",nil)
    router.ServeHTTP(w,req)
    if handlerID != "1" {
        t.Fatal("Routing failed with no Accept: header")
    }

    req.Header.Set("Accept","text/html")
    router.ServeHTTP(w,req)
    if handlerID != "1" {
        t.Fatal("Routing failed for explicit Accept: header")
    }

    req.Header.Set("Accept","text/plain")
    router.ServeHTTP(w,req)
    if handlerID != "unsupported" {
        t.Fatal("fallbackto Unsupported method failed")
    }
}
