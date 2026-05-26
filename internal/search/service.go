package search

import "strings"

type Service struct {
	bundles []TripBundle
}

func NewService() *Service {
	return &Service{
		bundles: seedBundles(),
	}
}

func (s *Service) Search(destination string) []TripBundle {
	if destination == "" {
		result := make([]TripBundle, len(s.bundles))
		copy(result, s.bundles)
		return result
	}

	lower := strings.ToLower(destination)
	parts := strings.FieldsFunc(lower, func(r rune) bool{
		return r==',' || r==' ' || r=='\t'
	})

	var filtered []TripBundle
	for _, bundle := range s.bundles {
		hotelLower := strings.ToLower(bundle.Hotel.Name) 
		arrLower := strings.ToLower(bundle.Flight.Arr)
		for _, part := range parts {
			if len(part) > 2 && (strings.Contains(hotelLower, part) ||
			strings.Contains(arrLower, part)) {
				filtered = append(filtered, bundle)
				break
			}
		}
	}
	return filtered
}

func seedBundles() []TripBundle {

	return []TripBundle{

		{

			ID: 1,

			Label: "Best value",

			Flight: Flight{
				Airline: "Lufthansa",
				Code:    "LH203",
				Dep:     "Kolkata",
				Arr:     "Berlin",
				Type:    "Economy",
			},

			Hotel: Hotel{
				Name:      "Hilton Berlin",
				Stars:     4,
				Breakfast: true,
			},

			Price:       487,
			CO2:         142,
			DurationMin: 165,
			InPolicy:    true,
			PolicyNote:  "",
		},

		{

			ID:    2,
			Label: "Premium comfort",
			Flight: Flight{
				Airline: "Emirates",
				Code:    "EK471",
				Dep:     "Kolkata",
				Arr:     "Berlin",
				Type:    "Business",
			},

			Hotel: Hotel{
				Name:      "Ritz-Carlton Berlin",
				Stars:     5,
				Breakfast: true,
			},
			Price:       1240,
			CO2:         210,
			DurationMin: 150,
			InPolicy:    false,
			PolicyNote:  "Exceeds budget cap of $800",
		},

		{
			ID:    3,
			Label: "Budget friendly",
			Flight: Flight{
				Airline: "Ryanair",
				Code:    "FR882",
				Dep:     "Kolkata",
				Arr:     "Berlin",
				Type:    "Economy",
			},
			Hotel: Hotel{
				Name:      "IBIS Berlin",
				Stars:     3,
				Breakfast: false,
			},
			Price:       299,
			CO2:         98,
			DurationMin: 180,
			InPolicy:    true,
			PolicyNote:  "",
		},

		{
			ID:    4,
			Label: "Quickest route",
			Flight: Flight{
				Airline: "British Airways",
				Code:    "BA144",
				Dep:     "Kolkata",
				Arr:     "London",
				Type:    "Economy",
			},
			Hotel: Hotel{
				Name:      "Holiday Inn London",
				Stars:     4,
				Breakfast: true,
			},
			Price:       620,
			CO2:         185,
			DurationMin: 120,
			InPolicy:    true,
			PolicyNote:  "",
		},
		{
			ID:    5,
			Label: "Luxury stay",
			Flight: Flight{
				Airline: "Singapore Airlines",
				Code:    "SQ325",
				Dep:     "Kolkata",
				Arr:     "Singapore",
				Type:    "Business",
			},
			Hotel: Hotel{
				Name:      "Marina Bay Sands",
				Stars:     5,
				Breakfast: true,
			},
			Price:       1890,
			CO2:         320,
			DurationMin: 240,
			InPolicy:    false,
			PolicyNote:  "Premium route requires VP approval",
		},
	}

}
