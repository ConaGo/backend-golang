package server

import (
	"net/http"
	"strconv"
	"strings"
)
const (
	layoutISO = "2006-01-02"
)
func StartServer(){
	http.HandleFunc("/", Serve)
	http.ListenAndServe(":8080", nil)
}

func Serve(w http.ResponseWriter, r *http.Request) {
/* 	//https://golang.org/pkg/time/#example_Tick
	t := time.Tick(30 * time.Minute)
	for next := range t {
		fmt.Printf("%v %s\n", next, statusUpdate())
	} */
	var h http.Handler
/* 	http.HandleFunc("/conferences", HandleConferences)
	http.HandleFunc("/questions", HandleQuestions) */
	
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