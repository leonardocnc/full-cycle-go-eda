package web

import (
	"balances/internal/usecase/balance/create_balance"
	"balances/internal/usecase/balance/get_balance"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type WebBalanceHandler struct {
	CreateBalanceUseCase create_balance.CreateBalanceUseCase
	GetBalanceUseCase    get_balance.GetBalanceUseCase
}

func NewWebBalanceHandler(
	c create_balance.CreateBalanceUseCase,
	g get_balance.GetBalanceUseCase,
) *WebBalanceHandler {
	return &WebBalanceHandler{
		CreateBalanceUseCase: c,
		GetBalanceUseCase:    g,
	}
}

func (h *WebBalanceHandler) CreateBalance(w http.ResponseWriter, r *http.Request) {
	var dto create_balance.CreateBalanceInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println("CreateBalance json.NewDecoder")
		fmt.Println("Error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	output, err := h.CreateBalanceUseCase.Execute(ctx, dto)
	if err != nil {
		fmt.Println("CreateBalance CreateBalanceUseCase.Execute")
		fmt.Println("Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		fmt.Println("CreateBalance json.NewEncoder(w).Encode")
		fmt.Println("Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *WebBalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	var dto get_balance.GetBalanceInputDTO

	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 {
		fmt.Println("GetBalance: Invalid URL path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	accountID := pathSegments[2]

	if accountID == "" {
		fmt.Println("GetBalance: AccountID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dto.AccountID = accountID

	output, err := h.GetBalanceUseCase.Execute(dto)
	if err != nil {
		fmt.Println("GetBalance Execute")
		fmt.Println("Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(output); err != nil {
		fmt.Println("GetBalance json.NewEncoder")
		fmt.Println("Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
