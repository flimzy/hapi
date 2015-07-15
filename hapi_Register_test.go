package hapi

import (
    "net/http"
    "testing"

    "github.com/gorilla/context"
    "github.com/julienschmidt/httprouter"   /* HTTP router */
)

func (h *HypermediaAPI) DoRegisterTest(name, requestedType, expectedType, expectedID string, t *testing.T) {
    negType,typeHandler := h.TypeAndHandler("GET","/",requestedType)
    if negType != expectedType {
        t.Fatalf("%s: TypeAndHandler returned '%s', expected '%s'\n", name, negType, expectedType)
    }
    w := new(mockResponseWriter)
    r,_ := http.NewRequest("GET","/",nil)

    if typeHandler != nil {
        typeHandler( w, r )
    }
    if id,ok := context.GetOk(r,"id"); ! ok {
        if ( len(expectedID) > 0 ) {
            t.Fatalf("%s: Error reading id after request\n", name)
        }
    } else if id != expectedID {
        t.Fatalf("%s: Handler identified itself as '%s', expected '%s'\n", name, id, expectedID)
    }
}

func TestRegister1(t *testing.T) {
    router := New()
    router.TestRegister("text/html","1")
    router.DoRegisterTest("text/html 1","text/html","text/html","1",t)
    router.DoRegisterTest("text/html 2","text/*","text/html","1",t)
    router.DoRegisterTest("text/html 3","*/*","text/html","1",t)
    router.DoRegisterTest("text/html 4","text/plain","","",t)
}

func TestHandle(t *testing.T) {
    var negotiatedType string
    var foo string
    handler := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
        negotiatedType = context.Get(r,"Content-Type").(string)
        foo = p.ByName("foo")
    }
    
    router := New()
    router.Handle("GET","/:foo","text/html",handler)
    w := new(mockResponseWriter)
    r,_ := http.NewRequest("GET","/bar",nil)
    router.ServeHTTP(w,r)
    if negotiatedType != "text/html" {
        t.Fatalf("httprouter.Handle set Content-Type to '%s', expected 'text/html'. httprouter context not working.\n", negotiatedType)
    }
    if foo != "bar" {
        t.Fatalf("httprouter.Handle set foo to '%s', expected 'bar'. httprouter.Params handling not working.\n", foo)
    }
}
