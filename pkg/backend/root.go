package backend

import (
	"html/template"
	"log"
	"net/http"

	"github.com/pojntfx/networkmate/pkg/models"
	"github.com/pojntfx/networkmate/pkg/persisters"
	"github.com/pojntfx/networkmate/pkg/templates"
)

type Backend struct {
	tpl       *template.Template
	persister *persisters.Persister
}

func NewBackend(persister *persisters.Persister) *Backend {
	return &Backend{
		persister: persister,
	}
}

func (b *Backend) Init() error {
	tpl, err := template.ParseFS(templates.FS, "*.html")
	if err != nil {
		return err
	}

	b.tpl = tpl

	return nil
}

type indexData struct {
	Contacts []models.Contact
}

func (b *Backend) HandleIndex(w http.ResponseWriter, r *http.Request) {
	contacts, err := b.persister.GetContacts(r.Context())
	if err != nil {
		log.Println("Could not fetch contacts:", err)

		http.Error(w, "Could not fetch contacts", http.StatusInternalServerError)

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "index.html", indexData{
		Contacts: contacts,
	}); err != nil {
		log.Println("Could not render template, continuing:", err)

		http.Error(w, "Could not render template", http.StatusInternalServerError)

		return
	}
}
