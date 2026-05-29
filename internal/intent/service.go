package intent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const geminiAPIBase = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"

type Service struct {
	geminiAPIKey string
	httpClient   *http.Client
}

func NewService(apiKey string) *Service {
	return &Service{
		geminiAPIKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *Service) ExtractIntent(ctx context.Context, text string) (*IntentResponse, error) {
	result, err := s.callGeminiAPI(ctx, text)
	if err != nil {
		result = s.fallbackExtract(text)
	}
	return result, nil
}

func (s *Service) callGeminiAPI(ctx context.Context, text string) (*IntentResponse, error) {
	prompt := fmt.Sprintf(`
You are a travel intent extraction engine.
Extract structured travel information from the user input.
STRICT INSTRUCTIONS:
* Return ONLY valid JSON
* Do NOT return markdown
* Do NOT explain anything
* Do NOT add extra fields
* If a value is missing, return empty string ""
* Extract source city after the word "from"
* Extract destination city after the word "to"
* Source and destination format must be:
 "City, Country"
Examples:
* "Bangalore, India"
* "Berlin, Germany"
* "Paris, France"
DATE RULES:
* Detect natural travel dates
* If only one weekday is mentioned:
 * dateFrom = upcoming weekday date
 * dateTo = next day
PURPOSE RULES:
* If text contains "client":
 purpose = "Client Meeting"
* If text contains "conference":
 purpose = "Conference"
* If text contains "training":
 purpose = "Training"
* Otherwise:
 purpose = "Business Travel"
Today's date is:
%s
And Today is :
%s
Return EXACTLY this JSON format:
{
"source":"",
"destination":"",
"dateFrom":"",
"dateTo":"",
"purpose":""
}
USER INPUT:
"%s"
`, time.Now().Local().Format("2006-01-02"), time.Now().Local().Weekday(), text)

	log.Print(prompt)

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Gemini request: %v", err)
	}

	apiUrl := fmt.Sprintf("%s?key=%s", geminiAPIBase, s.geminiAPIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl, bytes.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini API request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Gemini API request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API returned non-200 status: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Gemini API response: %v", err)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(responseBody, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Gemini API response: %v", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned by Gemini API")
	}

	textResponse := geminiResp.Candidates[0].Content.Parts[0].Text

	cleaned := s.cleanJSON(textResponse)
	log.Print("Cleaned")

	var result IntentResponse
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cleaned Gemini response: %v", err)
	}

	return &result, nil
}

func (s *Service) cleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)
	fmt.Println("Raw response:", raw)
	return raw
}

func (s *Service) fallbackExtract(text string) *IntentResponse {
	result := &IntentResponse{}
	destinations := map[string]string{
		"berlin":    "Berlin, Germany",
		"london":    "London, UK",
		"new york":  "New York, USA",
		"paris":     "Paris, France",
		"tokyo":     "Tokyo, Japan",
		"dubai":     "Dubai, UAE",
		"mumbai":    "Mumbai, India",
		"delhi":     "Delhi, India",
		"bangkok":   "Bangkok, Thailand",
		"singapore": "Singapore, Singapore",
	}
	lower := strings.ToLower(text)
	for key, val := range destinations {
		if strings.Contains(lower, key) {
			result.Destination = val
			break
		}
	}
	datePattern := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	dates := datePattern.FindAllString(text, -1)
	if len(dates) > 0 {
		result.DateFrom = dates[0]
		if len(dates) > 1 {
			result.DateTo = dates[1]
		} else {
			parsed, err := time.Parse("2006-01-02", dates[0])
			if err == nil {
				result.DateTo = parsed.AddDate(0, 0, 1).Format("2006-01-02")
			}
		}
	} else {
		today := time.Now()
		weekdays := map[string]time.Weekday{
			"monday":    time.Monday,
			"tuesday":   time.Tuesday,
			"wednesday": time.Wednesday,
			"thursday":  time.Thursday,
			"friday":    time.Friday,
			"saturday":  time.Saturday,
			"sunday":    time.Sunday,
		}
		for name, day := range weekdays {
			if strings.Contains(lower, name) {
				next := today
				diff := (day - today.Weekday() + 7) % 7
				if diff == 0 {
					diff = 7
				}
				next = today.AddDate(0, 0, int(diff))
				result.DateFrom = next.Format("2006-01-02")
				result.DateTo = next.AddDate(0, 0, 1).Format("2006-01-02")
				break
			}
		}
		if result.DateFrom == "" {
			if strings.Contains(lower, "tomorrow") {
				tomorrow := today.AddDate(0, 0, 1)
				result.DateFrom = tomorrow.Format("2006-01-02")
				result.DateTo = tomorrow.AddDate(0, 0, 1).Format("2006-01-02")
			} else {
				result.DateFrom = today.Format("2006-01-02")
				result.DateTo = today.AddDate(0, 0, 1).Format("2006-01-02")
			}
		}
	}
	purposes := []struct {
		keywords []string
		purpose  string
	}{
		{[]string{"client meeting", "meeting", "client"}, "Client Meeting"},
		{[]string{"conference", "summit", "convention"}, "Conference"},
		{[]string{"training", "workshop", "seminar"}, "Training"},
		{[]string{"interview", "recruitment", "hiring"}, "Interview"},
		{[]string{"audit", "inspection", "review"}, "Audit"},
		{[]string{"sales", "prospect", "pitch"}, "Sales Meeting"},
	}
	for _, p := range purposes {
		for _, kw := range p.keywords {
			if strings.Contains(lower, kw) {
				result.Purpose = p.purpose
				break
			}
		}
		if result.Purpose != "" {
			break
		}
	}
	if result.Purpose == "" {
		result.Purpose = "Business Travel"
	}
	return result
}
