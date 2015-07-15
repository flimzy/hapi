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
    "github.com/gorilla/context"
)

type HypermediaAPI struct {
    Router                  *httprouter.Router
    registeredTypes         []string
    typeHandlers            map[string]map[string]http.HandlerFunc
    // Configurable hapi.Handler which is called when the requested media
    // type (via Accept: header) cannot be served by a registered handler.
    // If it is not set, http.Error with http.StatusUnsupportedMediaType is used.
    UnsupportedMediaType    http.HandlerFunc
}

func New() *HypermediaAPI {
    return &HypermediaAPI{
        httprouter.New(),
        make([]string,0,1),
        make(map[string]map[string]http.HandlerFunc),
        defaultUnsupportedMediaType,
    }
}

func defaultUnsupportedMediaType (w http.ResponseWriter, r *http.Request) {
    http.Error(w,"415 Unsupported Media Type", http.StatusUnsupportedMediaType)
}

func (h *HypermediaAPI) GETAll(path string, handler http.HandlerFunc) {
    h.Register("GET",path,"*/*",handler)
}

func (h *HypermediaAPI) GET(path, ctype string, handler http.HandlerFunc) {
    h.Register("GET",path, ctype, handler)
}

func (h *HypermediaAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.Router.ServeHTTP(w,r)
}

// Handle() is an adaptor which allows the usage of an httprouter.Handle as a request handle
func (h *HypermediaAPI) Handle(method, path, ctype string, handler httprouter.Handle) {
    h.Register(method, path, ctype, func(w http.ResponseWriter, r *http.Request) {
        p := context.Get(r,"Params")
        handler(w,r,p.(httprouter.Params))
    })
}

// Handler() is an adapter which allows the usage of an http.Handler as a request handle
func (h *HypermediaAPI) Handler(method, path, ctype string, handler http.Handler) {
    h.Register(method, path, ctype, func(w http.ResponseWriter, r *http.Request) {
        handler.ServeHTTP(w, r)
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
func (h *HypermediaAPI) Register(method, path, ctype string, handler http.HandlerFunc) {
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
            accept := r.Header.Get("Accept")
            if len(accept) == 0 {
                accept = "*/*"
            }
            negotiatedType, typeHandler := h.TypeAndHandler(method,path,accept)
            if len(negotiatedType) == 0 {
                // Fall back to unsupported type
                typeHandler = h.UnsupportedMediaType
            }
//            w.Header.Set("Content-Type", negotiatedType)
            context.Set(r,"Params", p)
            context.Set(r,"Content-Type",negotiatedType)
            typeHandler( w, r )
            fmt.Fprintf(w, "Content-Type: %s\nparams: %s\n", context.Get(r,"Content-Type"),context.Get(r,"Params"))
            context.Clear(r)
            return
        }
        h.Router.Handle(method, path, wrapper)
        h.typeHandlers[key] = make(map[string]http.HandlerFunc)
    }
    h.registeredTypes = append(h.registeredTypes,ctypes...)
    for _,t := range ctypes {
        if strings.ContainsRune(t,'*') {
            panic(fmt.Sprintf("'%s' is not a valid media type; wildcards are not permitted.",t))
        }
        h.typeHandlers[key][t] = handler
    }
}

func (h *HypermediaAPI) TypeAndHandler(method, path, acceptHeader string) (negotiatedType string, typeHandler http.HandlerFunc) {
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
