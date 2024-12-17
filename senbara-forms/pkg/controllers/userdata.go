package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/pojntfx/senbara/senbara-forms/pkg/models"
)

const (
	exportTableNameJournalEntry = "journalEntry"
	exportTableNameContact      = "contact"
	exportTableNameDebt         = "debt"
	exportTableNameActivity     = "activity"
)

type userDataRow struct {
	TableName string
}

func (b *Controller) HandleUserData(w http.ResponseWriter, r *http.Request) {
	redirected, userData, status, err := b.authorize(w, r)
	if err != nil {
		log.Println(err)

		http.Error(w, err.Error(), status)

		return
	} else if redirected {
		return
	}

	w.Header().Set("Content-Type", "application/jsonl")
	w.Header().Set("Content-Disposition", `attachment; filename="senbara-forms-userdata.jsonl"`)

	encoder := json.NewEncoder(w)

	if err := b.persister.GetUserData(
		r.Context(),

		userData.Email,

		func(journalEntry models.GetJournalEntriesExportForNamespaceRow) error {
			journalEntry.TableName = exportTableNameJournalEntry

			if err := encoder.Encode(journalEntry); err != nil {
				return errors.Join(errCouldNotWriteResponse, err)
			}

			return nil
		},
		func(contact models.GetContactsExportForNamespaceRow) error {
			contact.TableName = exportTableNameContact

			if err := encoder.Encode(contact); err != nil {
				return errors.Join(errCouldNotWriteResponse, err)
			}

			return nil
		},
		func(debt models.GetDebtsExportForNamespaceRow) error {
			debt.TableName = exportTableNameDebt

			if err := encoder.Encode(debt); err != nil {
				return errors.Join(errCouldNotWriteResponse, err)
			}

			return nil
		},
		func(activity models.GetActivitiesExportForNamespaceRow) error {
			activity.TableName = exportTableNameActivity

			if err := encoder.Encode(activity); err != nil {
				return errors.Join(errCouldNotWriteResponse, err)
			}

			return nil
		},
	); err != nil {
		log.Println(errCouldNotFetchFromDB, err)

		http.Error(w, errCouldNotFetchFromDB.Error(), http.StatusInternalServerError)

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

	if err := b.persister.DeleteUserData(r.Context(), userData.Email); err != nil {
		log.Println(errCouldNotDeleteFromDB, err)

		http.Error(w, errCouldNotDeleteFromDB.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, userData.LogoutURL, http.StatusFound)
}
