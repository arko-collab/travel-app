package approval

type SubmitRequest struct {
	Destination string `json:"destination" validate:"required,min=2,max=255"`
	DateFrom    string `json:"dateFrom" validate:"required"`
	DateTo      string `json:"dateTo" validate:"required"`
	Purpose     string `json:"purpose" validate:"required,min=5,max=255"`
	FlightInfo  string `json:"flightInfo" validate:"max=1000"`
	HotelInfo   string `json:"hotelInfo" validate:"max=1000"`
	TotalCost   int    `json:"totalCost" validate:"required,min=0"`
	Notes       string `json:"notes" validate:"max=1000"`
}

type SubmitResponse struct {
	ApprovalID string `json:"approvalId"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
}

type EmailRequest struct {
	To          string
	ApprovalID  string
	Destination string
	DateFrom    string
	DateTo      string
	Purpose     string
	TotalCost   int
	Notes       string
	FlightInfo  string
	HotelInfo   string
}
