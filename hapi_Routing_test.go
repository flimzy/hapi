package hapi

import (
    "testing"
    "net/http"
//     "github.com/julienschmidt/httprouter"
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
    router.Register("GET","/","text/html", func (c *Context) { handlerID = "1" })
    router.UnsupportedMediaType = func (c *Context) { handlerID = "unsupported" }
    w := new(mockResponseWriter)
    req,_ := http.NewRequest("GET","/",nil)
    req.Header.Set("Accept","text/html")
    router.ServeHTTP(w,req)
    if handlerID != "1" {
        t.Fatal("routing failed")
    }

    req.Header.Set("Accept","text/plain")
    router.ServeHTTP(w,req)
    if handlerID != "unsupported" {
        t.Fatal("fallbackto Unsupported method failed")
    }
}
