[![Build Status](https://travis-ci.org/flimzy/hapi.svg?branch=master)](https://travis-ci.org/flimzy/hapi)

# What is HAPI?

HAPI, short for **Hypermedia API**, is a router and minimalistic framework for creating Hypermedia APIs (aka REST) for the Go language.
It is very much a work in progress, and in the very early stages of development. Therefore I do not recommend that you use it in anything
like a production environment.  If you are interested in partcipating in the development, your contribution is appreciated.

# Design goals

My primary design goal is to extend the traditional HTTP router to also handle multiple representations easily.

Most web frameworks and HTTP routers (upon which frameworks are often written) tend to focus on a single dimension of the HTTP request:
The URL.

An example from [gorilla mux](http://www.gorillatoolkit.org/pkg/mux):

    func main() {
        r := mux.NewRouter()
        r.HandleFunc("/", HomeHandler)
        r.HandleFunc("/products", ProductsHandler)
        r.HandleFunc("/articles", ArticlesHandler)
        http.Handle("/", r)
    }

For simple web sites, this can be enough (or nearly enough).  But in recent times, many routers have started to emphasise a second essential
dimension of the HTTP request: The HTTP method.

An example from [httprouter](https://github.com/julienschmidt/httprouter#web-frameworks-based-on-httprouter):

    func main() {
        // Initialize a router as usual
        router := httprouter.New()
        router.GET("/", Index)
        router.GET("/hello/:name", Hello)
        
This is obviously a big improvement for authors of REST APIs.  But it is my opinion that this doesn't go far enough, if you want to fully leverage
a Hypermedia API. The missing dimension is the media type.

Whenever your httprouter ignores a dimension you care about, you end up doing your own "routing" inside of your handler function. To handle GET and POST differently with gorilla/mux, for instance, you must check the request method within your handler, then call the appropriate function yourself.

Similarly, if you have one function that produces an HTML representation of a resource, and another that produces a PDF, your own handler must determine which media type is requested (by parsing the Accept: header in the request), then doing your own internal routing yourself.

By handling this in the router, it also makes it easier to provide a 406 error in the case that the client requests a type you can't provide, even if you only provide a single media type in your API.

We have similar problems if we need to consume requests with bodies of different media types. Perhaps you want to accept input from your clients as either JSON+Collection or as XML. You have to do all that bookkeeping and piping yourself.

Thus the primary goal of HAPI is to simplify the media type dimension of Hypermedia API routing.

## Compatibility

HAPI is a wrapper around [httprouter](https://github.com/julienschmidt/httprouter#web-frameworks-based-on-httprouter), and like httprouter, it provides a couple of functions to make it easy to use http.Handler, http.HandlerFunc, and even httprouter.Handle: hapi.Handler(), hapi.HandlerFunc(), and hapi.Handle() respectively.