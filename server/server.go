package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)
const (
	layoutISO = "2006-01-02"
)
var server *http.Server
var httpServerExitDone *sync.WaitGroup

func StartServer(){
    log.Printf("main: starting HTTP server")

    httpServerExitDone = &sync.WaitGroup{}

    httpServerExitDone.Add(1)
    server = startHttpServer(httpServerExitDone)	
}
func StopServer() {
    log.Printf("main: stopping HTTP server")
    if server == nil {
        log.Printf("no active server found")
        return
    }
    if err := server.Shutdown(context.TODO()); err != nil {
        panic(err) 
    }
    // wait for goroutine started in startHttpServer() to stop
    httpServerExitDone.Wait()

    log.Printf("main: done. exiting")
}
func startHttpServer(wg *sync.WaitGroup) *http.Server {
    srv := &http.Server{Addr: ":8080"}

    http.HandleFunc("/", Serve)

    go func() {
        defer wg.Done() // let main know we are done cleaning up

        // always returns error. ErrServerClosed on graceful close
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            // unexpected error. port in use?
            log.Fatalf("ListenAndServe(): %v", err)
        }
    }()

    // returning reference so caller can call Shutdown()
    return srv
}
func Serve(w http.ResponseWriter, r *http.Request) {
/* 	//https://golang.org/pkg/time/#example_Tick
	t := time.Tick(30 * time.Minute)
	for next := range t {
		fmt.Printf("%v %s\n", next, statusUpdate())
	} */
	var h http.Handler
	
	p := r.URL.Path
    switch {
    case match(p, "/conferences"):
        h = get(HandleConferences)
    case match(p, "/questions"):
        h = get(HandleQuestions)
    case match(p, "/questions/token"):
        h = get(HandleTokens)
    default:
        http.NotFound(w, r)
        return
    }
    h.ServeHTTP(w, r)
}
// get takes a HandlerFunc and wraps it to only allow the GET method
func get(h http.HandlerFunc) http.HandlerFunc {
    return allowMethod(h, "GET")
}
// allowMethod takes a HandlerFunc and wraps it in a handler that only
// responds if the request method is the given method, otherwise it
// responds with HTTP 405 Method Not Allowed.
func allowMethod(h http.HandlerFunc, method string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if method != r.Method {
            w.Header().Set("Allow", method)
            http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
            return
        }
        h(w, r)
    }
}
func match(path, pattern string, vars ...interface{}) bool {
    for ; pattern != "" && path != ""; pattern = pattern[1:] {
        switch pattern[0] {
        case '+':
            // '+' matches till next slash in path
            slash := strings.IndexByte(path, '/')
            if slash < 0 {
                slash = len(path)
            }
            segment := path[:slash]
            path = path[slash:]
            switch p := vars[0].(type) {
            case *string:
                *p = segment
            case *int:
                n, err := strconv.Atoi(segment)
                if err != nil || n < 0 {
                    return false
                }
                *p = n
            default:
                panic("vars must be *string or *int")
            }
            vars = vars[1:]
        case path[0]:
            // non-'+' pattern byte must match path byte
            path = path[1:]
        default:
            return false
        }
    }
    return path == "" && pattern == ""
}