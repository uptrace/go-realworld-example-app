package httputil

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

type ContextHandler struct {
	Ctx  context.Context
	Next http.Handler
}

func (h ContextHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.Next.ServeHTTP(w, req.WithContext(h.Ctx))
}

//------------------------------------------------------------------------------

type PanicHandler struct {
	Next http.Handler
}

func (h PanicHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<20)
			n := runtime.Stack(buf, true)
			fmt.Fprintf(os.Stderr, "panic: %v\n\n%s", err, buf[:n])
			os.Exit(1)
		}
	}()
	h.Next.ServeHTTP(w, req)
}
