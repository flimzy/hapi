// hapi provides a Hypermedia API (aka "true" REST) micro-framework/toolkit
package hapi

import (
    "fmt"
    "net/http"
    "github.com/julienschmidt/httprouter"   /* HTTP router */
//     "bitbucket.org/ww/goautoneg"            /* To parse Accept: headers */
)

type HypermediaAPI struct {
    Router  *httprouter.Router
}

type Context struct {
    Request *http.Request
    Params  httprouter.Params
    Stash   map[string]interface{}
}

type Handle func(*Context)

func New() *HypermediaAPI {
    return &HypermediaAPI{
        httprouter.New(),
    }
}

func (h *HypermediaAPI) GET(path string, handle Handle) {
    context := &Context{}
    wrapper := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
        context.Request = r
        context.Params = p
        context.Stash = make(map[string]interface{})
        handle( context )
        fmt.Fprintf(w, "%v", context)
    }
    h.Router.GET(path,wrapper)
}

func (h *HypermediaAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.Router.ServeHTTP(w,r)
    return
}
