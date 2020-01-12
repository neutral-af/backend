package airports

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestSearchCity(t *testing.T) {
	cases := []struct {
		query string
		IATA  string
	}{
		{"Toronto", "YYZ"},
		{"Montreal", "YUL"},
		{"Vancouver", "YVR"},
		{"Marrakesh", "RAK"},
	}

	for _, c := range cases {
		s := Search(c.query)
		fmt.Println(s)
		assert.NotZero(t, len(s))
		assert.Equal(t, c.IATA, s[0].IATA)
	}
}

func TestSearchAirportName(t *testing.T) {
	cases := []struct {
		query string
		IATA  string
	}{
		{"Pearson", "YYZ"},
		{"Trudeau", "YUL"},
		{"Gatwick", "LGW"},
		{"Heathrow", "LHR"},
		// {"Tegel", "TXL"},
	}

	for _, c := range cases {
		s := Search(c.query)
		assert.NotZero(t, len(s))
		assert.Equal(t, c.IATA, s[0].IATA)
	}
}

func TestSearchIATA(t *testing.T) {
	cases := []struct {
		IATA string
	}{
		{"YYZ"},
		{"YUL"},
		{"LGW"},
		{"LHR"},
		{"TXL"},
	}

	for _, c := range cases {
		s := Search(c.IATA)
		assert.NotZero(t, len(s))
		assert.Equal(t, c.IATA, s[0].IATA)
	}
}
