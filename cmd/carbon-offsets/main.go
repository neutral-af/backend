package main

import (
	"fmt"
	"log"

	"github.com/jasongwartz/carbon-offset-backend/lib/cloverly"
	"github.com/jasongwartz/carbon-offset-backend/lib/distance"
	"github.com/jasongwartz/carbon-offset-backend/lib/emissions"
)

func main() {
	distance := distance.TwoAirports("YYZ", "LHR")
	emissions := emissions.FlightCarbon(distance)
	x, err := cloverly.Estimate(emissions)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", x)
}
