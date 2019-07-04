package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jasongwartz/carbon-offset-backend/lib/cloverly"
	"github.com/jasongwartz/carbon-offset-backend/lib/distance"
	"github.com/jasongwartz/carbon-offset-backend/lib/emissions"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	distance := distance.TwoAirports(params.Get("departure"), params.Get("arrival"))
	emissions := emissions.FlightCarbon(distance)
	estimate, err := cloverly.Estimate(emissions)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estimate)
}
