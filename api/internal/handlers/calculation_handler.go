package handlers

import (
    "calculator/api/internal/service"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
)

type CalculationHandler struct {
    service *service.CalculationService
}

func NewCalculationHandler(service *service.CalculationService) *CalculationHandler {
    return &CalculationHandler{service: service}
}

type calculateRequest struct {
    Operand1  float64 `json:"operand1"`
    Operand2  float64 `json:"operand2"`
    Operation string  `json:"operation"`
}

func (h *CalculationHandler) Calculate(w http.ResponseWriter, r *http.Request) {
    var req calculateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    calc, err := h.service.Perform(req.Operand1, req.Operand2, req.Operation)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(calc); err != nil {
    	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    	return
	}
}

func (h *CalculationHandler) GetAllCalculations(w http.ResponseWriter, r *http.Request) {
    calcs, err := h.service.GetAll()
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(calcs); err != nil {
    	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    	return
	}
}

func (h *CalculationHandler) DeleteCalculation(w http.ResponseWriter, r *http.Request) {
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 4 {
        http.Error(w, "Invalid path", http.StatusBadRequest)
        return
    }
    id, err := strconv.Atoi(pathParts[3])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    if err := h.service.Delete(id); err != nil {
        http.Error(w, "Delete failed", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}