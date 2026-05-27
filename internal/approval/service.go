package approval

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/travel-api/build/internal/db"
	"gorm.io/gorm"
)

type Service struct {
	db             *gorm.DB
	sendGridAPIKey string
	fromEmail      string
	approverEmail  string
}

func NewService(database *gorm.DB, sendGridAPIKey, fromEmail, approverEmail string) *Service {
	return &Service{
		db:             database,
		sendGridAPIKey: sendGridAPIKey,
		fromEmail:      fromEmail,
		approverEmail:  approverEmail,
	}
}

func (s *Service) CreateApprovalRequest(req *SubmitRequest, ctx context.Context) (*SubmitResponse, error) {
	// Parse dates
	dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
	if err != nil {
		return nil, fmt.Errorf("invalid dateFrom format, expected YYYY-MM-DD: %w", err)
	}

	dateTo, err := time.Parse("2006-01-02", req.DateTo)
	if err != nil {
		return nil, fmt.Errorf("invalid dateTo format, expected YYYY-MM-DD: %w", err)
	}

	// Generate unique approval ID
	approvalID := fmt.Sprintf("APR-%d", time.Now().UnixNano())

	// Create approval request in database
	approvalReq := db.ApprovalRequest{
		ApprovalID:  approvalID,
		Destination: req.Destination,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		Purpose:     req.Purpose,
		FlightInfo:  req.FlightInfo,
		HotelInfo:   req.HotelInfo,
		TotalCost:   req.TotalCost,
		Notes:       req.Notes,
		Status:      "PENDING",
	}

	if err := s.db.WithContext(ctx).Create(&approvalReq).Error; err != nil {
		return nil, fmt.Errorf("failed to create approval request: %w", err)
	}

	// Trigger email notification to approver
	reqCopy := *req
	go s.TriggerApprovalEmail(approvalID, &reqCopy)

	return &SubmitResponse{
		ApprovalID: approvalID,
		Status:     "PENDING",
		Message:    "Approval request submitted successfully",
	}, nil
}

// TriggerApprovalEmail prepares and sends the approval email notification
func (s *Service) TriggerApprovalEmail(approvalID string, req *SubmitRequest) {
	emailReq := EmailRequest{
		To:          s.approverEmail,
		ApprovalID:  approvalID,
		Destination: req.Destination,
		DateFrom:    req.DateFrom,
		DateTo:      req.DateTo,
		Purpose:     req.Purpose,
		TotalCost:   req.TotalCost,
		Notes:       req.Notes,
		FlightInfo:  req.FlightInfo,
		HotelInfo:   req.HotelInfo,
	}

	if err := s.SendApprovalEmail(emailReq); err != nil {
		log.Printf("Failed to send approval email: %v", err)
		// Don't fail the request if email fails, just log it
	} else {
		log.Printf("Approval email sent successfully for %s", approvalID)
	}
}

// SendApprovalEmail sends a beautifully formatted HTML email to the approver
func (s *Service) SendApprovalEmail(req EmailRequest) error {
	from := mail.NewEmail("Travel Approval System", s.fromEmail)
	to := mail.NewEmail("Approver", req.To)
	subject := fmt.Sprintf("New Travel Approval Request - %s", req.ApprovalID)

	htmlContent := s.generateEmailHTML(req)
	plainTextContent := s.generateEmailPlainText(req)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.sendGridAPIKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid returned error status: %d", response.StatusCode)
	}

	log.Printf("Email sent successfully to %s for approval %s", req.To, req.ApprovalID)
	return nil
}

// generateEmailHTML creates a beautiful HTML email template
func (s *Service) generateEmailHTML(req EmailRequest) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
            font-weight: 600;
        }
        .content {
            padding: 30px;
        }
        .approval-id {
            background-color: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 15px;
            margin: 20px 0;
            font-size: 18px;
            font-weight: bold;
            color: #667eea;
        }
        .info-section {
            margin: 25px 0;
        }
        .info-row {
            display: flex;
            padding: 12px 0;
            border-bottom: 1px solid #eee;
        }
        .info-row:last-child {
            border-bottom: none;
        }
        .info-label {
            font-weight: 600;
            color: #555;
            min-width: 140px;
            display: flex;
            align-items: center;
        }
        .info-value {
            color: #333;
            flex: 1;
        }
        .icon {
            margin-right: 8px;
        }
        .highlight {
            background-color: #fff3cd;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            border-left: 4px solid #ffc107;
        }
        .cost-section {
            background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%);
            color: white;
            padding: 20px;
            border-radius: 8px;
            text-align: center;
            margin: 20px 0;
        }
        .cost-label {
            font-size: 14px;
            opacity: 0.9;
            margin-bottom: 5px;
        }
        .cost-amount {
            font-size: 32px;
            font-weight: bold;
        }
        .notes-section {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .notes-section h3 {
            margin-top: 0;
            color: #667eea;
            font-size: 16px;
        }
        .action-section {
            text-align: center;
            padding: 30px 0 10px 0;
        }
        .action-button {
            display: inline-block;
            padding: 14px 35px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            text-decoration: none;
            border-radius: 25px;
            font-weight: 600;
            margin: 0 10px;
            transition: transform 0.2s;
        }
        .footer {
            text-align: center;
            padding: 20px;
            color: #888;
            font-size: 12px;
            background-color: #f8f9fa;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✈️ New Travel Approval Request</h1>
        </div>
        
        <div class="content">
            <p>Hello,</p>
            <p>A new travel approval request has been submitted and requires your attention.</p>
            
            <div class="approval-id">
                🎫 Approval ID: %s
            </div>
            
            <div class="info-section">
                <div class="info-row">
                    <div class="info-label">
                        <span class="icon">🌍</span> Destination:
                    </div>
                    <div class="info-value">%s</div>
                </div>
                
                <div class="info-row">
                    <div class="info-label">
                        <span class="icon">📅</span> Travel Dates:
                    </div>
                    <div class="info-value">%s to %s</div>
                </div>
                
                <div class="info-row">
                    <div class="info-label">
                        <span class="icon">🎯</span> Purpose:
                    </div>
                    <div class="info-value">%s</div>
                </div>
            </div>
            
            <div class="cost-section">
                <div class="cost-label">Total Estimated Cost</div>
                <div class="cost-amount">$%d</div>
            </div>
            
            %s
            
            %s
            
            %s
            
            <div class="action-section">
                <p style="margin-bottom: 20px;">Please review and take action on this request:</p>
            </div>
        </div>
        
        <div class="footer">
            <p>This is an automated email from the Travel Approval System.</p>
            <p>Approval ID: %s | Generated on %s</p>
        </div>
    </div>
</body>
</html>
`,
		req.ApprovalID,
		req.Destination,
		req.DateFrom,
		req.DateTo,
		req.Purpose,
		req.TotalCost,
		s.formatFlightInfo(req.FlightInfo),
		s.formatHotelInfo(req.HotelInfo),
		s.formatNotes(req.Notes),
		req.ApprovalID,
		time.Now().Format("January 2, 2006 at 3:04 PM"),
	)
}

// formatFlightInfo formats flight information for HTML display
func (s *Service) formatFlightInfo(flightInfo string) string {
	if flightInfo == "" {
		return ""
	}
	return fmt.Sprintf(`
            <div class="notes-section">
                <h3>✈️ Flight Information</h3>
                <p>%s</p>
            </div>
	`, flightInfo)
}

// formatHotelInfo formats hotel information for HTML display
func (s *Service) formatHotelInfo(hotelInfo string) string {
	if hotelInfo == "" {
		return ""
	}
	return fmt.Sprintf(`
            <div class="notes-section">
                <h3>🏨 Hotel Information</h3>
                <p>%s</p>
            </div>
	`, hotelInfo)
}

// formatNotes formats additional notes for HTML display
func (s *Service) formatNotes(notes string) string {
	if notes == "" {
		return ""
	}
	return fmt.Sprintf(`
            <div class="highlight">
                <strong>📝 Additional Notes:</strong><br>
                %s
            </div>
	`, notes)
}

// generateEmailPlainText creates a plain text version of the email
func (s *Service) generateEmailPlainText(req EmailRequest) string {
	plainText := fmt.Sprintf(`
NEW TRAVEL APPROVAL REQUEST

Approval ID: %s

TRAVEL DETAILS:
Destination: %s
Travel Dates: %s to %s
Purpose: %s
Total Cost: $%d
`,
		req.ApprovalID,
		req.Destination,
		req.DateFrom,
		req.DateTo,
		req.Purpose,
		req.TotalCost,
	)

	if req.FlightInfo != "" {
		plainText += fmt.Sprintf("\nFlight Information:\n%s\n", req.FlightInfo)
	}

	if req.HotelInfo != "" {
		plainText += fmt.Sprintf("\nHotel Information:\n%s\n", req.HotelInfo)
	}

	if req.Notes != "" {
		plainText += fmt.Sprintf("\nAdditional Notes:\n%s\n", req.Notes)
	}

	plainText += "\n---\nThis is an automated email from the Travel Approval System."

	return plainText
}
