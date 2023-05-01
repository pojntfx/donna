package backend

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pojntfx/donna/internal/models"
	"github.com/pojntfx/donna/internal/templates"
	"github.com/pojntfx/donna/pkg/persisters"
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

type Backend struct {
	tpl       *template.Template
	persister *persisters.Persister

	oidcIssuer       string
	oidcClientID     string
	oidcClientSecret string
	oidcRedirectURL  string

	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
}

func NewBackend(
	persister *persisters.Persister,

	oidcIssuer,
	oidcClientID,
	oidcClientSecret string,
	oidcRedirectURL string,
) *Backend {
	return &Backend{
		persister: persister,

		oidcIssuer:       oidcIssuer,
		oidcClientID:     oidcClientID,
		oidcClientSecret: oidcClientSecret,
		oidcRedirectURL:  oidcRedirectURL,
	}
}

func (b *Backend) Init(ctx context.Context) error {
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
		ClientID:     b.oidcClientID,
		ClientSecret: b.oidcClientSecret,
		RedirectURL:  b.oidcRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "email", "email_verified"},
	}

	b.verifier = provider.Verifier(&oidc.Config{
		ClientID: b.oidcClientID,
	})

	return nil
}

type pageData struct {
	Page string
}

func (b *Backend) authenticate(w http.ResponseWriter, r *http.Request) (bool, string, error) {
	// refreshToken, err := r.Cookie(refreshTokenKey)
	// if err != nil {
	// 	if errors.Is(err, http.ErrNoCookie) {
	// 		http.Redirect(w, r, b.config.AuthCodeURL(b.oidcRedirectURL), http.StatusFound)

	// 		return true, "", nil
	// 	}

	// 	return false, "", err
	// }

	idToken, err := r.Cookie(idTokenKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// TODO: First try to get a new access token with refresh token and set it on the client. Only if that fails, redirect

			http.Redirect(w, r, b.config.AuthCodeURL(b.oidcRedirectURL), http.StatusFound)

			return true, "", nil
		}

		return false, "", err
	}

	id, err := b.verifier.Verify(r.Context(), idToken.Value)
	if err != nil {
		// TODO: First try to get a new access token with refresh token and set it on the client. Only if that fails, redirect

		http.Redirect(w, r, b.config.AuthCodeURL(b.oidcRedirectURL), http.StatusFound)

		return true, "", nil
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := id.Claims(&claims); err != nil {
		return false, "", err
	}

	if !claims.EmailVerified {
		return false, "", errEmailNotVerified
	}

	return false, claims.Email, nil
}

func (b *Backend) HandleIndex(w http.ResponseWriter, r *http.Request) {
	redirected, _, err := b.authenticate(w, r)
	if err != nil {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	}

	if redirected {
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)

		if err := b.tpl.ExecuteTemplate(w, "404.html", pageData{
			Page: "🕳️ Page not found",
		}); err != nil {
			log.Println(errCouldNotRenderTemplate, err)

			http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "index.html", pageData{
		Page: "🏠 Home",
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Backend) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	oauth2Token, err := b.config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenKey,
		Value:    oauth2Token.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	idToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     idTokenKey,
		Value:    idToken,
		Expires:  oauth2Token.Expiry,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	if err := b.tpl.ExecuteTemplate(w, "signing_in.html", pageData{
		Page: "🔓 Signing You In ...",
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

type journalEntryData struct {
	pageData
	Entry models.JournalEntry
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
			Page: "📓 Journal",
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
		Page: "➕ Add Journal Entry",
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

	rrating := r.FormValue("rating")
	if strings.TrimSpace(rrating) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	rating, err := strconv.Atoi(rrating)
	if err != nil {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	id, err := b.persister.CreateJournalEntry(r.Context(), title, body, int32(rating))
	if err != nil {
		log.Println(errCouldNotInsertIntoDB, err)

		http.Error(w, errCouldNotInsertIntoDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/journal/view?id=%v", id), http.StatusFound)
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

func (b *Backend) HandleEditJournal(w http.ResponseWriter, r *http.Request) {
	rid := r.FormValue("id")
	if strings.TrimSpace(rid) == "" {
		log.Println(errInvalidQueryParam)

		http.Error(w, errInvalidQueryParam.Error(), http.StatusUnprocessableEntity)

		return
	}

	id, err := strconv.Atoi(rid)
	if err != nil {
		log.Println(errInvalidQueryParam)

		http.Error(w, errInvalidQueryParam.Error(), http.StatusUnprocessableEntity)

		return
	}

	journalEntry, err := b.persister.GetJournalEntry(r.Context(), int32(id))
	if err != nil {
		log.Println(errCouldNotFetchFromDB, err)

		http.Error(w, errCouldNotFetchFromDB.Error(), http.StatusInternalServerError)

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "journal_edit.html", journalEntryData{
		pageData: pageData{
			Page: "✏️ Edit Journal Entry",
		},
		Entry: journalEntry,
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Backend) HandleUpdateJournal(w http.ResponseWriter, r *http.Request) {
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

	rrating := r.FormValue("rating")
	if strings.TrimSpace(rrating) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	rating, err := strconv.Atoi(rrating)
	if err != nil {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	if err := b.persister.UpdateJournalEntry(r.Context(), int32(id), title, body, int32(rating)); err != nil {
		log.Println(errCouldNotUpdateInDB, err)

		http.Error(w, errCouldNotInsertIntoDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, "/journal/view?id="+rid, http.StatusFound)
}

func (b *Backend) HandleViewJournal(w http.ResponseWriter, r *http.Request) {
	rid := r.FormValue("id")
	if strings.TrimSpace(rid) == "" {
		log.Println(errInvalidQueryParam)

		http.Error(w, errInvalidQueryParam.Error(), http.StatusUnprocessableEntity)

		return
	}

	id, err := strconv.Atoi(rid)
	if err != nil {
		log.Println(errInvalidQueryParam)

		http.Error(w, errInvalidQueryParam.Error(), http.StatusUnprocessableEntity)

		return
	}

	journalEntry, err := b.persister.GetJournalEntry(r.Context(), int32(id))
	if err != nil {
		log.Println(errCouldNotFetchFromDB, err)

		http.Error(w, errCouldNotFetchFromDB.Error(), http.StatusInternalServerError)

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "journal_view.html", journalEntryData{
		pageData: pageData{
			Page: journalEntry.Title,
		},
		Entry: journalEntry,
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Backend) HandleImprint(w http.ResponseWriter, r *http.Request) {
	if err := b.tpl.ExecuteTemplate(w, "imprint.html", pageData{
		Page: "ℹ️ Imprint",
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}
