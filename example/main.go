package main

import (
	"fmt"
	"github.com/daniloanp/muxUrlGen"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleUrlsWithFunc(rtr *mux.Router, urls []string, handler func(w http.ResponseWriter, r *http.Request)) {
	for _, url := range urls {
		rtr.HandleFunc(url, handler)
	}
}

func echoVars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, r.URL.String(), "\n\n")
	for k, i := range mux.Vars(r) {
		fmt.Fprintln(w, k, ": ", i)
	}
}

func main() {
	var url string
	var urls []string

	rtr := mux.NewRouter()

	url = "/handl/v1/{v}/v2/{v'}"
	rtr.HandleFunc(url, echoVars)
	// The above code works by we cannot missing any var OR Even

	// Using "long notation"
	urls = muxUrlGen.GetUrlVarsPermutations("/handlLong/v1/{v:[0-9]+}/v2/{v'}", true)
	HandleUrlsWithFunc(rtr, urls, echoVars)

	// Using "long notation" and Optional vars
	urls = muxUrlGen.GetUrlVarsPermutations("/handlLongOptional/v1/{v}?/v2/{v'}?", true)
	HandleUrlsWithFunc(rtr, urls, echoVars)

	// Using "short notation"
	urls = muxUrlGen.GetUrlVarsPermutations("/handlShort/{v1}/{v2}", false)
	HandleUrlsWithFunc(rtr, urls, echoVars)

	// Using "short notation" and Optional vars
	urls = muxUrlGen.GetUrlVarsPermutations("/handlShortOptional/{v1}?/{v2}?", false)
	HandleUrlsWithFunc(rtr, urls, echoVars)

	http.Handle("/", rtr)
	fmt.Println("Server running at http://127.0.0.1:8080")

	http.ListenAndServe(":8080", nil)
}
