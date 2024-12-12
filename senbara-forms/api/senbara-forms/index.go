package senbaraForms

import (
	"net/http"
	"os"

	_ "github.com/lib/pq"

	senbaraForms "github.com/pojntfx/senbara/senbara-forms"
	"github.com/pojntfx/senbara/senbara-forms/pkg/controllers"
	"github.com/pojntfx/senbara/senbara-forms/pkg/persisters"
	"github.com/pojntfx/senbara/senbara-forms/pkg/static"
)

var (
	p *persisters.Persister
	c *controllers.Controller
)

func SenbaraFormsHandler(
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
	mux.HandleFunc("/debts/edit", c.HandleEditDebt)

	mux.HandleFunc("/debts/create", c.HandleCreateDebt)
	mux.HandleFunc("/debts/settle", c.HandleSettleDebt)
	mux.HandleFunc("/debts/update", c.HandleUpdateDebt)

	mux.HandleFunc("/activities/add", c.HandleAddActivity)
	mux.HandleFunc("/activities/view", c.HandleViewActivity)
	mux.HandleFunc("/activities/edit", c.HandleEditActivity)

	mux.HandleFunc("/activities/create", c.HandleCreateActivity)
	mux.HandleFunc("/activities/delete", c.HandleDeleteActivity)
	mux.HandleFunc("/activities/update", c.HandleUpdateActivity)

	mux.HandleFunc("/authorize", c.HandleAuthorize)

	mux.Handle("/code/", http.StripPrefix("/code/", http.FileServer(http.FS(senbaraForms.FS))))

	mux.HandleFunc("/", c.HandleIndex)

	mux.ServeHTTP(w, r)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = r.URL.Query().Get("path")

	if p == nil {
		p = persisters.NewPersister(os.Getenv("POSTGRES_URL"))

		if err := p.Init(); err != nil {
			panic(err)
		}
	}

	if c == nil {
		c = controllers.NewController(
			p,

			os.Getenv("OIDC_ISSUER"),
			os.Getenv("OIDC_CLIENT_ID"),
			os.Getenv("OIDC_REDIRECT_URL"),

			os.Getenv("PRIVACY_URL"),
			os.Getenv("IMPRINT_URL"),
		)

		if err := c.Init(r.Context()); err != nil {
			panic(err)
		}
	}

	SenbaraFormsHandler(w, r, c)
}
