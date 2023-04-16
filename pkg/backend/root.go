package backend

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/pojntfx/donna/internal/models"
	"github.com/pojntfx/donna/internal/templates"
	"github.com/pojntfx/donna/pkg/persisters"
)

var (
	errCouldNotRenderTemplate = errors.New("could not render template")
	errCouldNotFetchFromDB    = errors.New("could not fetch from DB")
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
	tpl, err := template.New("").Funcs(template.FuncMap{
		"TruncateText": func(text string, length int) string {
			if len(text) <= length {
				return text
			}

			return text[:length] + "…"
		},
	}).ParseFS(templates.FS, "*.html")
	if err != nil {
		return err
	}

	b.tpl = tpl

	return nil
}

type indexData struct {
	Page string
}

func (b *Backend) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)

		if err := b.tpl.ExecuteTemplate(w, "404.html", indexData{
			Page: "Page not found",
		}); err != nil {
			log.Println(errCouldNotRenderTemplate, err)

			http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "index.html", indexData{
		Page: "Home",
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

type journalData struct {
	indexData
	Entries []models.JournalEntry
}

func (b *Backend) HandleJournal(w http.ResponseWriter, r *http.Request) {
	journalEntries, err := b.persister.GetJournalEntries(r.Context())
	if err != nil {
		log.Println(errCouldNotFetchFromDB, err)

		http.Error(w, errCouldNotFetchFromDB.Error(), http.StatusInternalServerError)

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "journal.html", journalData{
		indexData: indexData{
			Page: "Journal",
		},
		Entries: journalEntries,
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}
