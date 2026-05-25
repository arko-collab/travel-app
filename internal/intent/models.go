package intent

type IntentRequest struct {
	Text string `json:"text"`
}

type IntentResponse struct {
	Destination string `json:"destination"`
	DateFrom    string `json:"dateFrom"`
	DateTo      string `json:"dateTo"`
	Purpose     string `json:"purpose"`
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}
