package api

import (
	"ChiliOverFlow/pkg/db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
)

var (
	m *db.Bundle
)

func SetupAPI(mod *db.Bundle) {

	m = mod

	router := httprouter.New()

	router.GET("/health", healthCheck)
	router.GET("/inventory", getInventory)
	router.GET("/inventory/available", getInventoryQuantity)

	corsEnabledHandler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8081", corsEnabledHandler))
}

func healthCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	setHeader(w)
	io.WriteString(w, `{"alive": true}`)
}

func getInventory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	drinks, err := m.GetInventory()
	encodeValue(drinks, err, w)
}

func getInventoryQuantity(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q, err := m.GetInventoryTotalVariety()
	encodeValue(q, err, w)
}


func encodeValue(val interface{}, err error, w http.ResponseWriter) {
	setHeader(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(val)
}

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
