package service

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"net/http"
	"typo3uz/database"
)

// TGHandler describes a tree greed handler
type TGHandler struct {
	Session *sql.DB
}

// ServeHTTP implements ServeHTTP function to fulfill http.Handler interface
func (h *TGHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var response []byte
	switch r.Method {
	case http.MethodGet:
		res, err := database.GetLanguages(h.Session)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(err)
		}
		response = res
	case http.MethodPost:
		res, err := database.UpdateLanguages(r, h.Session)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(err)
		}
		response = res
	}
	// set this to allow Ajax requests from other origins
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
	w.Header().Set("Content-type", "application/json")
	_, err := w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
	}
}
