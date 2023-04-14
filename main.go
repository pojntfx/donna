package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pojntfx/networkmate/pkg/backend"
	"github.com/pojntfx/networkmate/pkg/persisters"
)

func main() {
	laddr := flag.String("laddr", ":1337", "Listen address")
	dbaddr := flag.String("dbaddr", "host=localhost user=postgres dbname=networkmate sslmode=disable", "Database address")

	flag.Parse()

	p := persisters.NewPersister(*dbaddr)

	if err := p.Init(); err != nil {
		panic(err)
	}

	b := backend.NewBackend(p)

	if err := b.Init(); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", b.HandleIndex)

	log.Println("Listening on", *laddr)

	panic(http.ListenAndServe(*laddr, mux))
}
