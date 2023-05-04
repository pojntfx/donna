package controllers

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pojntfx/donna/pkg/persisters"
	"github.com/pojntfx/donna/pkg/templates"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/oauth2"
)

var (
	errCouldNotRenderTemplate = errors.New("could not render template")
	errCouldNotFetchFromDB    = errors.New("could not fetch from DB")
	errCouldNotParseForm      = errors.New("could not parse form")
	errInvalidForm            = errors.New("could not use invalid form")
	errCouldNotInsertIntoDB   = errors.New("could not insert into DB")
	errCouldNotDeleteFromDB   = errors.New("could not delete from DB")
	errCouldNotUpdateInDB     = errors.New("could not update in DB")
	errInvalidQueryParam      = errors.New("could not use invalid query parameter")
	errCouldNotLogin          = errors.New("could not login")
	errEmailNotVerified       = errors.New("email not verified")
)

const (
	idTokenKey      = "id_token"
	refreshTokenKey = "refresh_token"
)

type Controller struct {
	tpl       *template.Template
	persister *persisters.Persister

	oidcIssuer      string
	oidcClientID    string
	oidcRedirectURL string

	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
}

func NewController(
	persister *persisters.Persister,

	oidcIssuer,
	oidcClientID,
	oidcRedirectURL string,
) *Controller {
	return &Controller{
		persister: persister,

		oidcIssuer:      oidcIssuer,
		oidcClientID:    oidcClientID,
		oidcRedirectURL: oidcRedirectURL,
	}
}

func (b *Controller) Init(ctx context.Context) error {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)

	tpl, err := template.New("").Funcs(template.FuncMap{
		"TruncateText": func(text string, length int) string {
			if len(text) <= length {
				return text
			}

			return text[:length] + "…"
		},
		"RenderMarkdown": func(text string) template.HTML {
			var buf bytes.Buffer
			if err := md.Convert([]byte(text), &buf); err != nil {
				panic(err)
			}

			return template.HTML(buf.String())
		},
	}).ParseFS(templates.FS, "*.html")
	if err != nil {
		return err
	}

	b.tpl = tpl

	provider, err := oidc.NewProvider(ctx, b.oidcIssuer)
	if err != nil {
		return err
	}

	b.config = &oauth2.Config{
		ClientID:    b.oidcClientID,
		RedirectURL: b.oidcRedirectURL,
		Endpoint:    provider.Endpoint(),
		Scopes:      []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "email", "email_verified"},
	}

	b.verifier = provider.Verifier(&oidc.Config{
		ClientID: b.oidcClientID,
	})

	return nil
}

func (b *Controller) HandleIndex(w http.ResponseWriter, r *http.Request) {
	redirected, authorizationData, err := b.authorize(w, r)
	if err != nil {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	} else if redirected {
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)

		if err := b.tpl.ExecuteTemplate(w, "404.html", pageData{
			authorizationData: authorizationData,

			Page: "🕳️ Page not found",
		}); err != nil {
			log.Println(errCouldNotRenderTemplate, err)

			http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "index.html", pageData{
		authorizationData: authorizationData,

		Page: "🏠 Home",
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Controller) HandleImprint(w http.ResponseWriter, r *http.Request) {
	if err := b.tpl.ExecuteTemplate(w, "imprint.html", pageData{
		Page: "ℹ️ Imprint",
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}
