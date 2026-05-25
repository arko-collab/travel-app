package intent

import (
	"encoding/json"
	"net/http"

	"github.com/travel-api/build/internal/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) HandleIntent(w http.ResponseWriter, r *http.Request) {
	var req IntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w,http.StatusBadRequest, "Invalid Request Body")
		return
	}

	if req.Text == "" {
		middleware.RespondError(w, http.StatusBadRequest, "Text field is required")
		return
	}

	ctx := r.Context()
	intent, err := h.svc.ExtractIntent(ctx, req.Text)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to extract intent")
		return
	}
	middleware.RespondJSON(w, http.StatusOK, intent)
}