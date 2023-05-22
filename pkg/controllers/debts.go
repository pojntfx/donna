package controllers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func (b *Controller) HandleAddDebt(w http.ResponseWriter, r *http.Request) {
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

	contact, err := b.persister.GetContact(r.Context(), int32(id), authorizationData.Email)
	if err != nil {
		log.Println(errCouldNotFetchFromDB, err)

		http.Error(w, errCouldNotFetchFromDB.Error(), http.StatusInternalServerError)

		return
	}

	if err := b.tpl.ExecuteTemplate(w, "debts_add.html", contactData{
		pageData: pageData{
			authorizationData: authorizationData,

			Page: "💸 Add Debt",
		},
		Entry: contact,
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Controller) HandleCreateDebt(w http.ResponseWriter, r *http.Request) {
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

	ryouOwe := r.FormValue("you_owe")
	if strings.TrimSpace(ryouOwe) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	youOwe, err := strconv.Atoi(ryouOwe)
	if err != nil {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	rcontactID := r.FormValue("contact_id")
	if strings.TrimSpace(rcontactID) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	contactID, err := strconv.Atoi(rcontactID)
	if err != nil {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	ramount := r.FormValue("amount")
	if strings.TrimSpace(rcontactID) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	amount, err := strconv.ParseFloat(ramount, 64)
	if err != nil {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	if youOwe == 1 {
		amount = -math.Abs(amount)
	} else {
		amount = math.Abs(amount)
	}

	currency := r.FormValue("currency")
	if strings.TrimSpace(currency) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	if _, err := b.persister.CreateDebt(
		r.Context(),

		amount,
		currency,

		int32(contactID),
		authorizationData.Email,
	); err != nil {
		log.Println(errCouldNotInsertIntoDB, err)

		http.Error(w, errCouldNotInsertIntoDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/contacts/view?id=%v", contactID), http.StatusFound)
}

func (b *Controller) HandleSettleDebt(w http.ResponseWriter, r *http.Request) {
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

	rcontactID := r.FormValue("contact_id")
	if strings.TrimSpace(rcontactID) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	contactID, err := strconv.Atoi(rcontactID)
	if err != nil {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	if err := b.persister.SettleDebt(
		r.Context(),

		int32(id),

		int32(contactID),
		authorizationData.Email,
	); err != nil {
		log.Println(errCouldNotUpdateInDB, err)

		http.Error(w, errCouldNotUpdateInDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/contacts/view?id=%v", contactID), http.StatusFound)
}
