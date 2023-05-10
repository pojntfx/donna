package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

	contactID := r.FormValue("contact_id")
	if strings.TrimSpace(contactID) == "" {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	id, err := strconv.Atoi(contactID)
	if err != nil {
		log.Println(errInvalidForm)

		http.Error(w, errInvalidForm.Error(), http.StatusUnprocessableEntity)

		return
	}

	belongsToNamespace, err := b.persister.ContactBelongsToNamespace(r.Context(), int32(id), authorizationData.Email)
	if err != nil {
		log.Println(errCouldNotFetchFromDB, err)

		http.Error(w, errCouldNotInsertIntoDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/debts/view?id=%v", id), http.StatusFound)
}
