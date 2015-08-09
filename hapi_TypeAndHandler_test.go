package hapi

import (
    "testing"
    "net/http"
)

func (h *HypermediaAPI) DoTypeAndHandlerTests(name, accept, expectedType, expectedID string, actualId *string, t *testing.T) {
    negType,negHandler := h.TypeAndHandler("GET","/",accept)
    if negType != expectedType {
        t.Fatalf("%s: TypeAndHandler returned '%s', expected '%s'\n", name, negType, expectedType)
    }
    *actualId = ""  // Reset the identifier for each test
    if negHandler != nil {
        w := new(mockResponseWriter)
        r,_ := http.NewRequest("GET","/bar",nil)
        p := Params{}
        negHandler(w,r,p)
    }
    if *actualId != expectedID {
        t.Fatalf("%s: Handler identified itself as '%s', expected '%s'\n", name, *actualId, expectedID)
    }
}

// Tests 'text/html' handler
func TestTypeNegotiator1(t *testing.T) {
    router := New()
    var id string
    router.Register("GET", "/", "text/html", func(w http.ResponseWriter, r *http.Request, p Params) { id = "1" })
    //                          Test Name       Accept:         ExpectedType    ExpectedID  ActualID
    router.DoTypeAndHandlerTests("text/html",   "text/html",    "text/html",    "1",        &id,    t)
    router.DoTypeAndHandlerTests("text/plain",  "text/plain",   "",             "",         &id,    t)
    router.DoTypeAndHandlerTests("text/*",      "text/*",       "text/html",    "1",        &id,    t)
    router.DoTypeAndHandlerTests("*/*",         "*/*",          "text/html",    "1",        &id,    t)
    router.DoTypeAndHandlerTests("foo/*",       "foo/*",        "",             "",         &id,    t)
    router.DoTypeAndHandlerTests("foo/*",       "",             "",             "",         &id,    t)
}

// Two routers: 1: text/html, 2: text/plain
func TestTypeNegotiator2(t *testing.T) {
    router := New()
    var id string
    router.Register("GET", "/", "text/html", func(w http.ResponseWriter, r *http.Request, p Params) { id = "1" })
    router.Register("GET", "/", "text/plain", func(w http.ResponseWriter, r *http.Request, p Params) { id = "2" })
    //                          Test Name       Accept:         ExpectedType    ExpectedID  ActualID
    router.DoTypeAndHandlerTests("text/html",   "text/html",    "text/html",    "1",        &id,    t)
    router.DoTypeAndHandlerTests("text/plain",  "text/plain",   "text/plain",   "2",        &id,    t)
    router.DoTypeAndHandlerTests("plain+html",  "text/plain, text/html",
                                                                "text/plain",   "2",        &id,    t)
    router.DoTypeAndHandlerTests("plain+html 2","text/plain;q=0.2, text/html",
                                                                "text/html",    "1",        &id,    t)
    router.DoTypeAndHandlerTests("text/*",      "text/*",       "text/html",    "1",        &id,    t)
    router.DoTypeAndHandlerTests("*/*",         "*/*",          "text/html",    "1",        &id,    t)
    router.DoTypeAndHandlerTests("foo/*",       "foo/*",        "",             "",         &id,    t)
}
