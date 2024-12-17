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

	file, _, err := r.FormFile("userData")
	if err != nil {
		log.Println(errCouldNotReadRequest, err)

		http.Error(w, errCouldNotReadRequest.Error(), http.StatusInternalServerError)

		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var rawRow json.RawMessage
		if err := decoder.Decode(&rawRow); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Println(errCouldNotReadRequest, err)

			http.Error(w, errCouldNotReadRequest.Error(), http.StatusInternalServerError)

			return
		}

		var userDataRow userDataRow
		if err := json.Unmarshal(rawRow, &userDataRow); err != nil {
			log.Println(errCouldNotReadRequest, err)

			http.Error(w, errCouldNotReadRequest.Error(), http.StatusInternalServerError)

			return
		}

		switch userDataRow.TableName {
		case exportTableNameJournalEntry:
			var journalEntry models.GetJournalEntriesExportForNamespaceRow
			if err := json.Unmarshal(rawRow, &journalEntry); err != nil {
				log.Println(errCouldNotReadRequest, err)

				http.Error(w, errCouldNotReadRequest.Error(), http.StatusInternalServerError)

				return
			}

			log.Println(journalEntry)

		case exportTableNameContact:
			var contact models.GetContactsExportForNamespaceRow
			if err := json.Unmarshal(rawRow, &contact); err != nil {
				log.Println(errCouldNotReadRequest, err)

				http.Error(w, errCouldNotReadRequest.Error(), http.StatusInternalServerError)

				return
			}

			log.Println(contact)

		case exportTableNameDebt:
			var debt models.GetDebtsExportForNamespaceRow
			if err := json.Unmarshal(rawRow, &debt); err != nil {
				log.Println(errCouldNotReadRequest, err)

				http.Error(w, errCouldNotReadRequest.Error(), http.StatusInternalServerError)

				return
			}

			log.Println(debt)

		case exportTableNameActivity:
			var activity models.GetActivitiesExportForNamespaceRow
			if err := json.Unmarshal(rawRow, &activity); err != nil {
				log.Println(errCouldNotReadRequest, err)

				http.Error(w, errCouldNotReadRequest.Error(), http.StatusInternalServerError)

				return
			}

			log.Println(activity)

		default:
			log.Println("Skipping import error:", errUnknownTableName, err)

			continue
		}
	}

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
