package backend

import (
	"html/template"
	"log"
	"net/http"

	"github.com/pojntfx/networkmate/pkg/templates"
)

type Backend struct {
	tpl *template.Template
}

func NewBackend() *Backend {
	return &Backend{}
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
	Name string
}

func (b *Backend) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if err := b.tpl.ExecuteTemplate(w, "index.html", indexData{
		Name: "Felicitas",
	}); err != nil {
		log.Println("Could not render template, continuing:", err)

		http.Error(w, "Could not render template", http.StatusInternalServerError)

		return
	}
}
