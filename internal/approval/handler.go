package approval

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/travel-api/build/internal/middleware"
)

type Handler struct {
	svc       *Service
	validator *validator.Validate
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc:       svc,
		validator: validator.New(),
	}
}

func (h *Handler) HandleApprovalRequest(w http.ResponseWriter, r *http.Request) {
	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ctx := r.Context();

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessage := formatValidationErrors(validationErrors)
		middleware.RespondError(w, http.StatusBadRequest, errorMessage)
		return
	}

	// Create approval request
	response, err := h.svc.CreateApprovalRequest(&req, ctx)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to create approval request")
		return
	}

	middleware.RespondJSON(w, http.StatusCreated, response)
}

// formatValidationErrors formats validator errors into readable messages
func formatValidationErrors(errors validator.ValidationErrors) string {
	errorMsg := "Validation failed: "
	for i, err := range errors {
		if i > 0 {
			errorMsg += ", "
		}
		switch err.Tag() {
		case "required":
			errorMsg += err.Field() + " is required"
		case "min":
			errorMsg += err.Field() + " must be at least " + err.Param() + " characters"
		case "max":
			errorMsg += err.Field() + " must be at most " + err.Param() + " characters"
		default:
			errorMsg += err.Field() + " is invalid"
		}
	}
	return errorMsg
}
