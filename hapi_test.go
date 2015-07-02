package hapi

import (
    "testing"
    "net/http"
//     "net/http/httptest"
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


func TestNew(t *testing.T) {
    h := New()
    if h.Router == nil {
        t.Fatal("h.Router not initilized")
    }
}

// Test a handler for 'text/html' only
func TestTypeNegotiator(t *testing.T) {
    handler := func(c *Context) {}
    handlersMap := make(map[string]Handle)
    handlersMap["text/html"] = handler

    // Accept: matches exactly
    negType,negHandler := TypeNegotiator("text/html",handlersMap)
    if negType != "text/html" {
        t.Fatalf("TypeNegotiator returned '%s', not '%s'\n", negType, "text/html")
    }
    if negHandler == nil {
        t.Fatal("No handler returned")
    }

    // Accept: Doesn't match
    negType,negHandler = TypeNegotiator("text/plain",handlersMap)
    if negType != "" {
        t.Fatalf("TypeNegotiator returned '%s', not '%s'\n", negType, "")
    }
    if negHandler != nil {
        t.Fatal("Handler returned inappropriately")
    }

    // Accept: doesn't exist
    negType,negHandler = TypeNegotiator("",handlersMap)
    if negType != "text/html" {
        t.Fatalf("TypeNegotiator returned '%s', not '%s'\n", negType, "text/html")
    }
    if negHandler == nil {
        t.Fatal("No handler returned")
    }

    // Accept: */*
    negType,negHandler = TypeNegotiator("*/*",handlersMap)
    if negType != "text/html" {
        t.Fatalf("TypeNegotiator returned '%s', not '%s'\n", negType, "text/html")
    }
    if negHandler == nil {
        t.Fatal("No handler returned")
    }

    // Accept: text/*
    negType,negHandler = TypeNegotiator("*/*",handlersMap)
    if negType != "text/html" {
        t.Fatalf("TypeNegotiator returned '%s', not '%s'\n", negType, "text/html")
    }
    if negHandler == nil {
        t.Fatal("No handler returned")
    }

    // Accept: foo/*
    negType,negHandler = TypeNegotiator("text/plain",handlersMap)
    if negType != "" {
        t.Fatalf("TypeNegotiator returned '%s', not '%s'\n", negType, "")
    }
    if negHandler != nil {
        t.Fatal("Handler returned inappropriately")
    }


}

// func TestGET(t *testing.T) {
//     router := New()
//     routed := false
//     router.GET("/foo","*/*",func(c *Context) {
//         routed = true
//         if c.NegotiatedType != "text/html" {
//             t.Fatal("Expected a negotiated type of 'text/html'")
//         }
//     })
//     
//     w := new(mockResponseWriter)
//     r, _ := http.NewRequest("GET", "/foo", nil)
//     router.ServeHTTP(w,r)
//     
//     if !routed {
//         t.Fatal("Request was not routed")
//     }
// }
