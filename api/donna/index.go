package donna

import (
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/pojntfx/donna/pkg/backend"
	"github.com/pojntfx/donna/pkg/persisters"
	"github.com/pojntfx/donna/pkg/static"
)

var (
	p *persisters.Persister
	b *backend.Backend
)

func DonnaHandler(
	w http.ResponseWriter,
	r *http.Request,
	b *backend.Backend,
) {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	mux.HandleFunc("/journal", b.HandleJournal)
	mux.HandleFunc("/journal/add", b.HandleAddJournal)
	mux.HandleFunc("/journal/edit", b.HandleEditJournal)
	mux.HandleFunc("/journal/view", b.HandleViewJournal)

	mux.HandleFunc("/journal/create", b.HandleCreateJournal)
	mux.HandleFunc("/journal/delete", b.HandleDeleteJournal)
	mux.HandleFunc("/journal/update", b.HandleUpdateJournal)

	mux.HandleFunc("/contacts", b.HandleContacts)
	mux.HandleFunc("/contacts/add", b.HandleAddContact)

	mux.HandleFunc("/contacts/create", b.HandleCreateContact)
	mux.HandleFunc("/contacts/delete", b.HandleDeleteContact)

	mux.HandleFunc("/imprint", b.HandleImprint)

	mux.HandleFunc("/authorize", b.HandleAuthorize)

	mux.HandleFunc("/", b.HandleIndex)

	mux.ServeHTTP(w, r)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = r.URL.Query().Get("path")

	if p == nil {
		p = persisters.NewPersister(os.Getenv("DATABASE_URL"))

		if err := p.Init(); err != nil {
			panic(err)
		}
	}

	if b == nil {
		b = backend.NewBackend(p, os.Getenv("OIDC_ISSUER"), os.Getenv("OIDC_CLIENT_ID"), os.Getenv("OIDC_REDIRECT_URL"))

		if err := b.Init(r.Context()); err != nil {
			panic(err)
		}
	}

	DonnaHandler(w, r, b)
}
