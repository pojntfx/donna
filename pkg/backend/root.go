package backend

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/pojntfx/donna/internal/models"
	"github.com/pojntfx/donna/internal/templates"
	"github.com/pojntfx/donna/pkg/persisters"
)

var (
	errCouldNotRenderTemplate = errors.New("could not render template")
	errCouldNotFetchFromDB    = errors.New("could not fetch from DB")
	errCouldNotParseForm      = errors.New("could not parse form")
	errInvalidForm            = errors.New("could not use invalid form")
	errCouldNotInsertIntoDB   = errors.New("could not insert into DB")
	errCouldNotDeleteFromDB   = errors.New("could not delete from DB")
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

type pageData struct {
	Page string
}

func (b *Backend) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)

		if err := b.tpl.ExecuteTemplate(w, "404.html", pageData{
			Page: "Page not found",
		}); err != nil {
			log.Println(errCouldNotRenderTemplate, err)

			http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "index.html", pageData{
		Page: "Home",
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

type journalData struct {
	pageData
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
		pageData: pageData{
			Page: "Journal",
		},
		Entries: journalEntries,
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Backend) HandleAddJournal(w http.ResponseWriter, r *http.Request) {
	if err := b.tpl.ExecuteTemplate(w, "journal_add.html", pageData{
		Page: "Journal",
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Backend) HandleCreateJournal(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(errCouldNotParseForm, err)

		http.Error(w, errCouldNotParseForm.Error(), http.StatusInternalServerError)

		return
	}

	title := r.FormValue("title")
	if strings.TrimSpace(title) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	body := r.FormValue("body")
	if strings.TrimSpace(body) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	if err := b.persister.CreateJournalEntry(r.Context(), title, body); err != nil {
		log.Println(errCouldNotInsertIntoDB, err)

		http.Error(w, errCouldNotInsertIntoDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, "/journal", http.StatusFound)
}

func (b *Backend) HandleDeleteJournal(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(errCouldNotParseForm, err)

		http.Error(w, errCouldNotParseForm.Error(), http.StatusInternalServerError)

		return
	}

	rid := r.FormValue("id")
	if strings.TrimSpace(rid) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	id, err := strconv.Atoi(rid)
	if err != nil {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	if err := b.persister.DeleteJournalEntry(r.Context(), int32(id)); err != nil {
		log.Println(errCouldNotDeleteFromDB, err)

		http.Error(w, errCouldNotDeleteFromDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, "/journal", http.StatusFound)
}
