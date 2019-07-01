package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jasongwartz/carbon-offset-backend/lib/cloverly"
	"github.com/jasongwartz/carbon-offset-backend/lib/distance"
	"github.com/jasongwartz/carbon-offset-backend/lib/emissions"
)

const port = "8080"

func main() {
	http.HandleFunc("/flight/estimate", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		distance := distance.TwoAirports(params.Get("departure"), params.Get("arrival"))
		emissions := emissions.FlightCarbon(distance)
		estimate, err := cloverly.Estimate(emissions)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(estimate)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
