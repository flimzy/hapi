package hapi

import (
    "testing"
//     "net/http"
)

func (h *HypermediaAPI) DoRegisterTest(name, requestedType, expectedType, expectedID string, t *testing.T) {
    negType,negHandler := h.TypeAndHandler("GET","/",requestedType)
    if negType != expectedType {
        t.Fatalf("%s: TypeAndHandler returned '%s', expected '%s'\n", name, negType, expectedType)
    }
    context := &Context{}
    context.Stash = make(map[string]interface{})
    if negHandler != nil {
        negHandler( context )
    }
    var id string
    if context.Stash["id"] != nil {
        id = context.Stash["id"].(string)
    }
    if id != expectedID {
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
