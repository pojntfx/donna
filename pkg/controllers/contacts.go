package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/pojntfx/donna/pkg/models"
)

type contactsData struct {
	pageData
	Entries []models.Contact
}

type contactData struct {
	pageData
	Entry models.Contact
}

func (b *Controller) HandleContacts(w http.ResponseWriter, r *http.Request) {
	redirected, authorizationData, err := b.authorize(w, r)
	if err != nil {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	} else if redirected {
		return
	}

	contacts, err := b.persister.GetContacts(r.Context(), authorizationData.Email)
	if err != nil {
		log.Println(errCouldNotFetchFromDB, err)

		http.Error(w, errCouldNotFetchFromDB.Error(), http.StatusInternalServerError)

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "contacts.html", contactsData{
		pageData: pageData{
			authorizationData: authorizationData,

			Page: "👥 Contacts",
		},
		Entries: contacts,
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Controller) HandleAddContact(w http.ResponseWriter, r *http.Request) {
	redirected, authorizationData, err := b.authorize(w, r)
	if err != nil {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	} else if redirected {
		return
	}

	if err := b.tpl.ExecuteTemplate(w, "contacts_add.html", pageData{
		authorizationData: authorizationData,

		Page: "🤝 Add Contact",
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Controller) HandleCreateContact(w http.ResponseWriter, r *http.Request) {
	redirected, authorizationData, err := b.authorize(w, r)
	if err != nil {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	} else if redirected {
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println(errCouldNotParseForm, err)

		http.Error(w, errCouldNotParseForm.Error(), http.StatusInternalServerError)

		return
	}

	firstName := r.FormValue("first_name")
	if strings.TrimSpace(firstName) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	lastName := r.FormValue("last_name")
	if strings.TrimSpace(lastName) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	email := r.FormValue("email")
	if _, err := mail.ParseAddress(email); err != nil {
		log.Println(err)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	nickname := r.FormValue("nickname")

	pronouns := r.FormValue("pronouns")
	if strings.TrimSpace(pronouns) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	id, err := b.persister.CreateContact(
		r.Context(),
		firstName,
		lastName,
		nickname,
		email,
		pronouns,
		authorizationData.Email,
	)
	if err != nil {
		log.Println(errCouldNotInsertIntoDB, err)

		http.Error(w, errCouldNotInsertIntoDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/contacts/view?id=%v", id), http.StatusFound)
}

func (b *Controller) HandleDeleteContact(w http.ResponseWriter, r *http.Request) {
	redirected, authorizationData, err := b.authorize(w, r)
	if err != nil {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	} else if redirected {
		return
	}

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

	if err := b.persister.DeleteContact(r.Context(), int32(id), authorizationData.Email); err != nil {
		log.Println(errCouldNotDeleteFromDB, err)

		http.Error(w, errCouldNotDeleteFromDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, "/contacts", http.StatusFound)
}

func (b *Controller) HandleViewContact(w http.ResponseWriter, r *http.Request) {
	redirected, authorizationData, err := b.authorize(w, r)
	if err != nil {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	} else if redirected {
		return
	}

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

	journalEntry, err := b.persister.GetContact(r.Context(), int32(id), authorizationData.Email)
	if err != nil {
		log.Println(errCouldNotFetchFromDB, err)

		http.Error(w, errCouldNotFetchFromDB.Error(), http.StatusInternalServerError)

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "contacts_view.html", contactData{
		pageData: pageData{
			authorizationData: authorizationData,

			Page: journalEntry.FirstName + " " + journalEntry.LastName,
		},
		Entry: journalEntry,
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}
