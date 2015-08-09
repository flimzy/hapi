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
//     "log"
    "strings"

    "net/http"
    "github.com/julienschmidt/httprouter"   /* HTTP router */
    "bitbucket.org/ww/goautoneg"            /* To parse Accept: headers */
)

// Exact copoy of Handle
type Handle func(http.ResponseWriter, *http.Request, Params)
type Params httprouter.Params
func (ps Params) ByName(name string) string {
    for i := range ps {
        if ps[i].Key == name {
            return ps[i].Value
        }
    }
    return ""
}

type HypermediaAPI struct {
    Router                  *httprouter.Router
    registeredTypes         []string
    typeHandlers            map[string]map[string]Handle
    // Configurable hapi.Handler which is called when the requested media
    // type (via Accept: header) cannot be served by a registered handler.
    // If it is not set, http.Error with http.StatusUnsupportedMediaType is used.
    UnsupportedMediaType    Handle
}

func New() *HypermediaAPI {
    return &HypermediaAPI{
        httprouter.New(),
        make([]string,0,1),
        make(map[string]map[string]Handle),
        defaultUnsupportedMediaType,
    }
}

func defaultUnsupportedMediaType (w http.ResponseWriter, r * http.Request, p Params) {
    http.Error(w,"415 Unsupported Media Type", http.StatusUnsupportedMediaType)
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

// Handler() is an adapter which allows the usage of an http.Handler as a request handle
func (h *HypermediaAPI) Handler(method, path, ctype string, handler http.Handler) {
    h.Register(method, path, ctype, func(w http.ResponseWriter, r *http.Request, p Params) {
        handler.ServeHTTP(w,r)
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
            h.dispatch(method,path,w,r,p)
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

func (h *HypermediaAPI) dispatch(method, path string,w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    accept := r.Header.Get("Accept")
    if len(accept) == 0 {
        accept = "*/*"
    }
    negotiatedType, typeHandler := h.TypeAndHandler(method,path,accept)
    if len(negotiatedType) == 0 {
        // Fall back to unsupported type
        typeHandler = h.UnsupportedMediaType
    } else {
        w.Header().Set("Content-Type",negotiatedType)
    }
    // Finagle the data type
    p := make(Params,0,len(params))
    p = append(p, params...)
    typeHandler(w,r,p)
    return
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
