package donna

import (
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/pojntfx/donna/pkg/controllers"
	"github.com/pojntfx/donna/pkg/persisters"
	"github.com/pojntfx/donna/pkg/static"
)

var (
	p *persisters.Persister
	c *controllers.Controller
)

func DonnaHandler(
	w http.ResponseWriter,
	r *http.Request,
	c *controllers.Controller,
) {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	mux.HandleFunc("/journal", c.HandleJournal)
	mux.HandleFunc("/journal/add", c.HandleAddJournal)
	mux.HandleFunc("/journal/edit", c.HandleEditJournal)
	mux.HandleFunc("/journal/view", c.HandleViewJournal)

	mux.HandleFunc("/journal/create", c.HandleCreateJournal)
	mux.HandleFunc("/journal/delete", c.HandleDeleteJournal)
	mux.HandleFunc("/journal/update", c.HandleUpdateJournal)

	mux.HandleFunc("/contacts", c.HandleContacts)
	mux.HandleFunc("/contacts/add", c.HandleAddContact)
	mux.HandleFunc("/contacts/edit", c.HandleEditContact)
	mux.HandleFunc("/contacts/view", c.HandleViewContact)

	mux.HandleFunc("/contacts/create", c.HandleCreateContact)
	mux.HandleFunc("/contacts/delete", c.HandleDeleteContact)
	mux.HandleFunc("/contacts/update", c.HandleUpdateContact)

	mux.HandleFunc("/debts/add", c.HandleAddDebt)

	mux.HandleFunc("/debts/create", c.HandleCreateDebt)
	mux.HandleFunc("/debts/settle", c.HandleSettleDebt)

	mux.HandleFunc("/imprint", c.HandleImprint)

	mux.HandleFunc("/authorize", c.HandleAuthorize)

	mux.HandleFunc("/", c.HandleIndex)

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

	if c == nil {
		c = controllers.NewController(p, os.Getenv("OIDC_ISSUER"), os.Getenv("OIDC_CLIENT_ID"), os.Getenv("OIDC_REDIRECT_URL"))

		if err := c.Init(r.Context()); err != nil {
			panic(err)
		}
	}

	DonnaHandler(w, r, c)
}
