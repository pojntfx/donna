package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (b *Controller) HandleUserData(w http.ResponseWriter, r *http.Request) {
	redirected, userData, status, err := b.authorize(w, r)
	if err != nil {
		log.Println(err)

		http.Error(w, err.Error(), status)

		return
	} else if redirected {
		return
	}

	allUserData, err := b.persister.GetAllUserData(r.Context(), userData.Email)
	if err != nil {
		log.Println(errCouldNotFetchFromDB, err)

		http.Error(w, errCouldNotFetchFromDB.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(allUserData); err != nil {
		log.Println(errCouldNotWriteResponse, err)

		http.Error(w, errCouldNotWriteResponse.Error(), http.StatusInternalServerError)

		return
	}
}

func (b *Controller) HandleCreateUserData(w http.ResponseWriter, r *http.Request) {
	redirected, _, status, err := b.authorize(w, r)
	if err != nil {
		log.Println(err)

		http.Error(w, err.Error(), status)

		return
	} else if redirected {
		return
	}

	// TODO: Import user data

	http.Redirect(w, r, "/contacts", http.StatusFound)
}

func (b *Controller) HandleDeleteUserData(w http.ResponseWriter, r *http.Request) {
	redirected, userData, status, err := b.authorize(w, r)
	if err != nil {
		log.Println(err)

		http.Error(w, err.Error(), status)

		return
	} else if redirected {
		return
	}

	if err := b.persister.DeleteAllUserData(r.Context(), userData.Email); err != nil {
		log.Println(errCouldNotDeleteFromDB, err)

		http.Error(w, errCouldNotDeleteFromDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, userData.LogoutURL, http.StatusFound)
}
