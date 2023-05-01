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
	"github.com/pojntfx/donna/api/donna"
	"github.com/pojntfx/donna/pkg/backend"
	"github.com/pojntfx/donna/pkg/persisters"
)

func main() {
	laddr := flag.String("laddr", ":1337", "Listen address (port can also be set with `PORT` env variable)")
	dbaddr := flag.String("dbaddr", "postgresql://postgres@localhost:5432/donna?sslmode=disable", "Database address (can also be set using `DATABASE_URL` env variable)")
	oidcIssuer := flag.String("oidc-issuer", "", "OIDC Issuer (i.e. https://pojntfx.eu.auth0.com/) (can also be set using the OIDC_ISSUER env variable)")
	oidcClientID := flag.String("oidc-client-id", "", "OIDC Client ID (i.e. myoidcclientid) (can also be set using the OIDC_CLIENT_ID env variable)")
	oidcRedirectURL := flag.String("oidc-redirect-url", "http://localhost:1337/authorize", "OIDC redirect URL (can also be set using the OIDC_REDIRECT_URL env variable)")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if v := os.Getenv("PORT"); v != "" {
		log.Println("Using port from PORT env variable")

		la, err := net.ResolveTCPAddr("tcp", *laddr)
		if err != nil {
			panic(err)
		}

		p, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}

		la.Port = p
		*laddr = la.String()
	}

	if v := os.Getenv("DATABASE_URL"); v != "" {
		log.Println("Using database address from DATABASE_URL env variable")

		*dbaddr = v
	}

	if v := os.Getenv("OIDC_ISSUER"); v != "" {
		log.Println("Using OIDC issuer from OIDC_ISSUER env variable")

		*oidcIssuer = v
	}

	if v := os.Getenv("OIDC_CLIENT_ID"); v != "" {
		log.Println("Using OIDC client ID from OIDC_CLIENT_ID env variable")

		*oidcClientID = v
	}

	if v := os.Getenv("OIDC_REDIRECT_URL"); v != "" {
		log.Println("Using OIDC redirect URL from OIDC_REDIRECT_URL env variable")

		*oidcRedirectURL = v
	}

	p := persisters.NewPersister(*dbaddr)

	if err := p.Init(); err != nil {
		panic(err)
	}

	b := backend.NewBackend(p, *oidcIssuer, *oidcClientID, *oidcRedirectURL)

	if err := b.Init(ctx); err != nil {
		panic(err)
	}

	log.Println("Listening on", *laddr)

	panic(http.ListenAndServe(*laddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		donna.DonnaHandler(w, r, b)
	})))
}
