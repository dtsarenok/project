package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/madnight/gocui-calculator/pkg/parser"
)

// Request структура для получения данных от клиента
type Request struct {
	Expression string `json:"expression"`
}

// Response структура для отправки результата
type Response struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	expr := strings.ReplaceAll(req.Expression, ",", ".")
	result, err := parser.Calculate(expr)
	if err != nil {
		log.Printf("Error calculating expression: %s\n", expr)
		response := &Response{Error: "Expression is not valid"}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonResponse)
		return
	}

	response := &Response{Result: result}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	fmt.Println("Server started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
