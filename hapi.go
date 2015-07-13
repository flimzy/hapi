// hapi provides a Hypermedia API (aka "true" REST) micro-framework/toolkit
package hapi

// hapi uses httprouter (github.com/julienschmidt/httprouter) for HTTP routing, and exposes
// all of the httprouter functionality through hapi.Router:

// func main() {
//     hapi := hapi.New()
//     hapi.GET(...) // hapi version
//     hapi.Router.GET(...) // Underlying httprouter version
// }

import (
    "fmt"
    "log"
    "strings"

    "net/http"
    "github.com/julienschmidt/httprouter"   /* HTTP router */
    "bitbucket.org/ww/goautoneg"            /* To parse Accept: headers */
)

type Handle func(*Context)

type HypermediaAPI struct {
    Router          *httprouter.Router
    registeredTypes []string
    typeHandlers    map[string]map[string]Handle
}

type Context struct {
    Writer              http.ResponseWriter
    Request             *http.Request
    Params              httprouter.Params
    NegotiatedType      string
    Stash               map[string]interface{}
}

func New() *HypermediaAPI {
    return &HypermediaAPI{
        httprouter.New(),
        make([]string,0,1),
        make(map[string]map[string]Handle),
    }
}

func (h *HypermediaAPI) GETAll(path string, handle Handle) {
    h.Register("GET",path,"*/*",handle)
}

func (h *HypermediaAPI) GET(path, ctype string, handle Handle) {
    h.Register("GET",path, ctype, handle)
}

func (h *HypermediaAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.Router.ServeHTTP(w,r)
}

// Handle() is an adaptor which allows the usage of an httprouter.Handle as a request handle
func (h *HypermediaAPI) Handle(method, path, ctype string, handle httprouter.Handle) {
    h.Register(method, path, ctype, func(c *Context) {
        handle(c.Writer,c.Request,c.Params)
    })
}

// Handler() is an adapter which allows the usage of an http.Handler as a request handle
func (h *HypermediaAPI) Handler(method, path, ctype string, handler http.Handler) {
    h.Register(method, path, ctype, func(c *Context) {
        handler.ServeHTTP(c.Writer, c.Request)
    })
}

// HandlerFunc() is an adaptor which allows the usage of an http.HandlerFunc as a request handle
func (h *HypermediaAPI) HandlerFunc(method, path, ctype string, handlerFunc http.HandlerFunc) {
    h.Handler(method, path, ctype, handlerFunc)
}

// Register() registers a handler method to handle a specific Method/path/content-type combination
// Method and path ought to be self-explanatory.
// The content type argument should be a space-separated list of valid content types. Wildcards
// are not permitted.
func (h *HypermediaAPI) Register(method, path, ctype string, handle Handle) {
    key := method + " " + path
    ctypes := strings.Split(ctype," ")
    if typeHandlers, registered := h.typeHandlers[key]; registered {
        for _,t := range ctypes {
            if _, ok := typeHandlers[t]; ok {
                panic(fmt.Sprintf("a handle is already registered for method %s, path '%s', type %s",method,path,ctype))
            }
        }
    } else {
        wrapper := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
            negotiatedType, typeHandler := h.TypeAndHandler(method,path,r.Header.Get("Accept"))
            if len(negotiatedType) == 0 {
                log.Printf("We can't serve the requested type(s): %s\n", r.Header.Get("Accept"))
                // Fall back to unsupported type
            }
//            w.Header.Set("Content-Type", negotiatedType)
            context := &Context{
                w,
                r,
                p,
                negotiatedType,
                make(map[string]interface{}),
            }
            typeHandler( context )
            fmt.Fprintf(w, "%v", context)
            return
        }
        h.Router.Handle(method, path, wrapper)
        h.typeHandlers[key] = make(map[string]Handle)
    }
    h.registeredTypes = append(h.registeredTypes,ctypes...)
    for _,t := range ctypes {
        if strings.ContainsRune(t,'*') {
            panic(fmt.Sprintf("'%s' is not a valid media type; wildcards are not permitted.",t))
        }
        h.typeHandlers[key][t] = handle
    }
}

func (h *HypermediaAPI) TypeAndHandler(method, path, acceptHeader string) (negotiatedType string, typeHandler Handle) {
    negotiatedType = goautoneg.Negotiate(acceptHeader,h.registeredTypes)
    if len(negotiatedType) == 0 {
        // This means we can't serve the requested type
        return
    }
    if handler,ok := h.typeHandlers[method + " " + path][negotiatedType]; ok {
        typeHandler = handler
        return
    }
//    typeHandler = nil
    return
}
