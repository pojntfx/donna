package backend

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/pojntfx/networkmate/internal/templates"
	"github.com/pojntfx/networkmate/pkg/persisters"
)

var (
	errCouldNotRenderTemplate = errors.New("could not render template")
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
