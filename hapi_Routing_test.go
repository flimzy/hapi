package hapi

import (
    "testing"
    "net/http"

//     "github.com/julienschmidt/httprouter"
)

func TestRouting(t *testing.T) {
    router := New()
    handlerID := ""
    router.Register("GET","/","text/html", func (_ http.ResponseWriter, _ *http.Request) { handlerID = "1" })
    router.UnsupportedMediaType = func (_ http.ResponseWriter, _ *http.Request) { handlerID = "unsupported" }
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
