package delivery

import (
	"encoding/json"
	"net/http"
	"payment-service/internal/usecase"

	"github.com/go-chi/chi"
)

type HttpHandler struct {
	uc *usecase.PaymentUsecase
}

func NewHttpHandler(uc *usecase.PaymentUsecase) *HttpHandler {
	return &HttpHandler{uc: uc}
}

func (h *HttpHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req usecase.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.uc.TransferFunds(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *HttpHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	refID := chi.URLParam(r, "refId")
	if refID == "" {
		respondWithError(w, http.StatusBadRequest, "reference ID is required")
		return
	}

	resp, err := h.uc.GetTransactionByRef(refID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "transaction not found")
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *HttpHandler) TopUp(w http.ResponseWriter, r *http.Request) {
	var req usecase.TopUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.uc.TopUpWallet(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
