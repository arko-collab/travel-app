package search

import (
	"encoding/json"
	"net/http"

	"github.com/travel-api/build/internal/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}
	results := h.svc.Search(req.Destination)
	if results == nil {
		results = []TripBundle{}
	}
	middleware.RespondJSON(w, http.StatusOK, results)
}
