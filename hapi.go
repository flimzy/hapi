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

func (h *HypermediaAPI) Register(method, path, ctype string, handle Handle) {
    key := method + " " + path
    ctypes := strings.Split(ctype," ")
    typeHandlers, registered := h.typeHandlers[key]
    for _,t := range ctypes {
        if _, ok := typeHandlers[t]; ok {
            panic(fmt.Sprintf("a handle is already registered for method %s, path '%s', type %s",method,path,ctype))
        }
    }
    
    if !registered {
        wrapper := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
            negotiatedType, typeHandler := TypeNegotiator(r.Header.Get("Accept"), h.typeHandlers[key])
log.Printf("Accept: %s\n", r.Header.Get("Accept"))
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
    for _,t := range ctypes {
        log.Printf("Registering for %s\n", t)
        h.typeHandlers[key][t] = handle
    }
}

func TypeNegotiator(acceptHeader string, typeHandlers map[string]Handle) (negotiatedType string, typeHandler Handle) {
    availableTypes := make([]string,len(typeHandlers))
    i := 0
    for k,_ := range typeHandlers {
        availableTypes[i] = k
        i++
    }
    if len(acceptHeader) == 0 {
        acceptHeader = "*/*"
    }
    negotiatedType = Negotiate(acceptHeader,availableTypes)
    if len(negotiatedType) == 0 {
        // This means we can't serve the requested type
        return
    }
    for _,negType := range []string{ negotiatedType, negotiatedType[0:strings.Index(negotiatedType,"/")]+"/*", "*/*" } {
        if handler,ok := typeHandlers[negType]; ok {
            typeHandler = handler
            return
        }
    }
    typeHandler = nil
    return
}

// Borrowed from goautoneg, and adapted
func Negotiate(header string, alternatives []string) (content_type string) {
    asp := make([][]string, 0, len(alternatives))
    for _, ctype := range alternatives {
        asp = append(asp, strings.SplitN(ctype, "/", 2))
    }
    for _, clause := range goautoneg.ParseAccept(header) {
        for i, ctsp := range asp {
            if clause.Type == ctsp[0] && clause.SubType == ctsp[1] {
                content_type = alternatives[i]
                return
            }
            if clause.Type == ctsp[0] && clause.SubType == "*" {
                content_type = alternatives[i]
                return
            }
            if clause.Type == "*" && clause.SubType == "*" {
                content_type = alternatives[i]
                return
            }
            if clause.Type == ctsp[0] && ctsp[1] == "*" {
                content_type = clause.Type + "/" + clause.SubType
                return
            }
            if ctsp[0] == "*" && ctsp[1] == "*" {
                content_type = clause.Type + "/" + clause.SubType
                return
            }
        }
    }
    return
}