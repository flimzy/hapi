package hapi

import (
    "net/http"
    "testing"
)

func (h *HypermediaAPI) DoRegisterTest(name, requestedType, expectedType, expectedID string, actualId *string, t *testing.T) {
    negType,negHandler := h.TypeAndHandler("GET","/",requestedType)
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

func TestRegister1(t *testing.T) {
    router := New()
    var id string
    router.Register("GET", "/", "text/html", func(w http.ResponseWriter, r *http.Request, p Params) { id = "1" })
    //                          Test Name       Accept:         ExpectedType    ExpectedID  ActualID
    router.DoRegisterTest("text/html 1",        "text/html",    "text/html",    "1",        &id,    t)
    router.DoRegisterTest("text/html 2",        "text/*",       "text/html",    "1",        &id,    t)
    router.DoRegisterTest("text/html 3",        "*/*",          "text/html",    "1",        &id,    t)
    router.DoRegisterTest("text/html 4",        "text/plain",   "",             "",         &id,    t)
}
