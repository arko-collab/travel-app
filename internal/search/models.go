package search

type SearchRequest struct {
	Destination string `json:"destination"`
}

type TripBundle struct {
	ID          int    `json:"id"`
	Label       string `json:"label"`
	Flight      Flight `json:"flight"`
	Hotel       Hotel  `json:"hotel"`
	Price       int    `json:"price"`
	CO2         int    `json:"co2"`
	DurationMin int    `json:"durationMin"`
	InPolicy    bool   `json:"inPolicy"`
	PolicyNote  string `json:"policyNote"`
}

type Hotel struct {
	Name      string `json:"name"`
	Stars     int    `json:"stars"`
	Breakfast bool   `json:"breakfast"`
}

type Flight struct {
	Airline string `json:"airline"`
	Code    string `json:"code"`
	Dep     string `json:"dep"`
	Arr     string `json:"arr"`
	Type    string `json:"type"`
}
