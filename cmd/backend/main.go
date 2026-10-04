package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	name := flag.String("name", "backend", "name returned in responses")
	listenAddr := flag.String("listen", ":9001", "address for the backend HTTP listener")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "served by %s\npath: %s\n", *name, r.URL.Path)
	})

	server := &http.Server{
		Addr:              *listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("backend %q listening on %s", *name, *listenAddr)
	log.Fatal(server.ListenAndServe())
}
