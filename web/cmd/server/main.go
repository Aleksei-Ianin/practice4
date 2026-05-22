package main

import (
	"calculator/web/internal/client"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var apiClient *client.APIClient

func main() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8081"
	}
	apiClient = client.NewAPIClient(apiURL)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/calculate", proxyCalculate)
	http.HandleFunc("/api/calculations", proxyGetAll)
	http.HandleFunc("/api/calculations/", proxyDelete)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Web server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func proxyCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Operand1  float64 `json:"operand1"`
		Operand2  float64 `json:"operand2"`
		Operation string  `json:"operation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	calc, err := apiClient.Calculate(req.Operand1, req.Operand2, req.Operation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calc)
}

func proxyGetAll(w http.ResponseWriter, r *http.Request) {
	calcs, err := apiClient.GetAllCalculations()
	if err != nil {
		http.Error(w, "API error", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calcs)
}

func proxyDelete(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	err = apiClient.DeleteCalculation(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}