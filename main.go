package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/pojntfx/networkmate/pkg/backend"
)

func main() {
	laddr := flag.String("laddr", ":1337", "Listen address")

	flag.Parse()

	b := backend.NewBackend()

	if err := b.Init(); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", b.HandleIndex)

	log.Println("Listening on", *laddr)

	panic(http.ListenAndServe(*laddr, mux))
}
