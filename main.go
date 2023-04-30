package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pojntfx/donna/internal/static"
	"github.com/pojntfx/donna/pkg/backend"
	"github.com/pojntfx/donna/pkg/persisters"
)

func main() {
	laddr := flag.String("laddr", ":1337", "Listen address (port can also be set with `PORT` env variable)")
	dbaddr := flag.String("dbaddr", "postgresql://postgres@localhost:5432/donna?sslmode=disable", "Database address (can also be set using `DATABASE_URL` env variable)")
	oidcIssuer := flag.String("oidc-issuer", "", "OIDC Issuer (i.e. https://pojntfx.eu.auth0.com/) (can also be set using the OIDC_ISSUER env variable)")
	oidcClientID := flag.String("oidc-client-id", "", "OIDC Client ID (i.e. myoidcclientid) (can also be set using the OIDC_CLIENT_ID env variable)")
	oidcRedirectURL := flag.String("oidc-redirect-url", "http://localhost:1337/login", "OIDC redirect URL (can also be set using the OIDC_REDIRECT_URL env variable)")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	if oi := os.Getenv("OIDC_ISSUER"); oi != "" {
		log.Println("Using OIDC issuer from OIDC_ISSUER env variable")

		*oidcIssuer = oi
	}

	if oc := os.Getenv("OIDC_CLIENT_ID"); oc != "" {
		log.Println("Using OIDC client ID from OIDC_CLIENT_ID env variable")

		*oidcIssuer = oc
	}

	p := persisters.NewPersister(*dbaddr)

	if err := p.Init(); err != nil {
		panic(err)
	}

	b := backend.NewBackend(p, *oidcIssuer, *oidcClientID, *oidcRedirectURL)

	if err := b.Init(ctx); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	mux.HandleFunc("/journal", b.HandleJournal)
	mux.HandleFunc("/journal/add", b.HandleAddJournal)
	mux.HandleFunc("/journal/edit", b.HandleEditJournal)
	mux.HandleFunc("/journal/view", b.HandleViewJournal)

	mux.HandleFunc("/journal/create", b.HandleCreateJournal)
	mux.HandleFunc("/journal/delete", b.HandleDeleteJournal)
	mux.HandleFunc("/journal/update", b.HandleUpdateJournal)

	mux.HandleFunc("/imprint", b.HandleImprint)

	mux.HandleFunc("/", b.HandleIndex)

	log.Println("Listening on", *laddr)

	panic(http.ListenAndServe(*laddr, mux))
}
