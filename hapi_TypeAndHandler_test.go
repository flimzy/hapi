package hapi

import (
    "net/http"
    "testing"

    "github.com/gorilla/context"
)

func (h *HypermediaAPI) DoTypeAndHandlerTests(name, accept, expectedType, expectedID string, t *testing.T) {
    negType,negHandler := h.TypeAndHandler("GET","/",accept)
    if negType != expectedType {
        t.Fatalf("%s: TypeAndHandler returned '%s', expected '%s'\n", name, negType, expectedType)
    }
    w := new(mockResponseWriter)
    r,_ := http.NewRequest("GET","/",nil)
    if negHandler != nil {
        negHandler( w, r )
    }
    if id,ok := context.GetOk(r,"id"); ! ok {
        if ( len(expectedID) > 0 ) {
            t.Fatalf("%s: Error reading id after request\n", name)
        }
    } else if id != expectedID {
        t.Fatalf("%s: Handler identified itself as '%s', expected '%s'\n", name, id, expectedID)
    }
}

// Tests 'text/html' handler
func TestTypeNegotiator1(t *testing.T) {
    router := New()
    router.TestRegister("text/html","1")
    //                          Test Name       Accept:         ExpectedType    ExpectedID
    router.DoTypeAndHandlerTests("text/html",   "text/html",    "text/html",    "1",    t)
    router.DoTypeAndHandlerTests("text/plain",  "text/plain",   "",             "",     t)
    router.DoTypeAndHandlerTests("text/*",      "text/*",       "text/html",    "1",    t)
    router.DoTypeAndHandlerTests("*/*",         "*/*",          "text/html",    "1",    t)
    router.DoTypeAndHandlerTests("foo/*",       "foo/*",        "",             "",     t)
}

// Two routers: 1: text/html, 2: text/plain
func TestTypeNegotiator2(t *testing.T) {
    router := New()
    router.TestRegister("text/html","1")
    router.TestRegister("text/plain","2")
    //                          Test Name       Accept:         ExpectedType    ExpectedID
    router.DoTypeAndHandlerTests("text/html",   "text/html",    "text/html",    "1",    t)
    router.DoTypeAndHandlerTests("text/plain",  "text/plain",   "text/plain",   "2",    t)
    router.DoTypeAndHandlerTests("plain+html",  "text/plain, text/html",
                                                                "text/plain",   "2",    t)
    router.DoTypeAndHandlerTests("plain+html 2","text/plain;q=0.2, text/html",
                                                                "text/html",    "1",    t)
    router.DoTypeAndHandlerTests("text/*",      "text/*",       "text/html",    "1",    t)
    router.DoTypeAndHandlerTests("*/*",         "*/*",          "text/html",    "1",    t)
    router.DoTypeAndHandlerTests("foo/*",       "foo/*",        "",             "",     t)
}
