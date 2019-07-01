package emissions

const (
	shortHaulMaxDistance   = 3700 // kilometres
	greatCircleFactor      = 1.09
	radiativeForcingFactor = 1.9
)

// FlightCarbon returns an estimate of the per-passenger
// carbon emissions for a flight of the given distance
// in units of kilograms CO2
func FlightCarbon(km float64) float64 {
	var emissions float64
	if km < shortHaulMaxDistance {
		emissions = 78.7
	} else {
		emissions = 103.1
	}
	return km * emissions * greatCircleFactor * radiativeForcingFactor / 1000
}
