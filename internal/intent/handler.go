package intent

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/travel-api/build/internal/middleware"
	"github.com/travel-api/build/internal/search"
)

type Handler struct {
	svc       *Service
	searchSvc *search.Service
}

func NewHandler(svc *Service, searchSvc *search.Service) *Handler {
	return &Handler{svc: svc, searchSvc: searchSvc}
}

func (h *Handler) HandleIntent(w http.ResponseWriter, r *http.Request) {
	var req IntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid Request Body")
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
	log.Printf("Extracted intent: %+v", intent)
	searchRes := h.searchSvc.Search(intent.Destination)
	if(searchRes == nil){
		searchRes = []search.TripBundle{}
	}
	middleware.RespondJSON(w, http.StatusOK, searchRes)
}
