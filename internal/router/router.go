package router

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Router is a traced version of httprouter.Router.
type Router struct {
	*httprouter.Router

	handler http.Handler
}

// New returns a new router augmented with tracing.
func New(opts ...otelhttp.Option) *Router {
	r := &Router{Router: httprouter.New()}

	const nInternalOpts = 1
	o := make([]otelhttp.Option, len(opts)+nInternalOpts)
	// Put this first so it can be overridden by the user.
	o[0] = otelhttp.WithSpanNameFormatter(r.name)
	if len(opts) > 0 {
		copy(o[1:], opts)
	}
	r.handler = otelhttp.NewHandler(r.Router, "", o...)

	return r
}

// ServeHTTP serves req writing a response to w and tracing the process.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}

// name is a SpanFormatter for the otelhttp instrumentation.
func (r *Router) name(_ string, req *http.Request) string {
	path := req.URL.Path
	_, params, trailing := r.Router.Lookup(req.Method, path)
	for _, param := range params {
		path = strings.Replace(path, param.Value, ":"+param.Key, 1)
	}
	if trailing {
		path = strings.TrimSuffix(path, "/")
	}

	return "HTTP " + req.Method + " " + path
}
