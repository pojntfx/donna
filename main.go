package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pojntfx/donna/pkg/backend"
	"github.com/pojntfx/donna/pkg/persisters"
)

func main() {
	laddr := flag.String("laddr", ":1337", "Listen address (port can also be set with `PORT` env variable)")
	dbaddr := flag.String("dbaddr", "postgresql://postgres@localhost:5432/donna?sslmode=disable", "Database address (can also be set using `DATABASE_URL` env variable)")

	flag.Parse()

	if p := os.Getenv("PORT"); p != "" {
		log.Println("Using port from PORT env variable")

		la, err := net.ResolveTCPAddr("tcp", *laddr)
		if err != nil {
			panic(err)
		}

		p, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}

		la.Port = p
		*laddr = la.String()
	}

	if da := os.Getenv("DATABASE_URL"); da != "" {
		log.Println("Using database address from DATABASE_URL env variable")

		*dbaddr = da
	}

	p := persisters.NewPersister(*dbaddr)

	if err := p.Init(); err != nil {
		panic(err)
	}

	b := backend.NewBackend(p)

	if err := b.Init(); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/journal/add", b.HandleAddJournal)
	mux.HandleFunc("/journal/create", b.HandleCreateJournal)
	mux.HandleFunc("/journal", b.HandleJournal)
	mux.HandleFunc("/", b.HandleIndex)

	log.Println("Listening on", *laddr)

	panic(http.ListenAndServe(*laddr, mux))
}
