package handler

import (
	"el-music-be/internal/database"
	"el-music-be/internal/middleware"
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentHandler struct {
	Store *database.PostgresStore
	Snap  snap.Client
}

func NewPaymentHandler(store *database.PostgresStore) *PaymentHandler {
	var s snap.Client
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	return &PaymentHandler{
		Store: store,
		Snap:  s,
	}
}

type ChargeRequest struct {
	Plan string `json:"plan"`
}

func (h *PaymentHandler) HandleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}

	user, err := h.Store.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var req ChargeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var amount int64
	var planName string
	if req.Plan == "monthly" {
		amount = 59000
		planName = "El Music Premium (Bulanan)"
	} else {
		http.Error(w, "Invalid plan", http.StatusBadRequest)
		return
	}

	orderID := "ELMUSIC-" + uuid.New().String()

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    req.Plan,
				Price: amount,
				Qty:   1,
				Name:  planName,
			},
		},
	}

	snapResp, err := h.Snap.CreateTransaction(snapReq)
	if err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"payment_url":    snapResp.RedirectURL,
		"transaction_id": snapResp.Token,
	})
}
