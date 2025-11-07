package get_transaction_history

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/samber/lo"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	adminModule *admin.Module
}

func NewHandler(adminModule *admin.Module) *Handler {
	return &Handler{
		adminModule: adminModule,
	}
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		studentIDStr := chi.URLParam(r, "student_id")
		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		from, to, err := req.ToTime()
		if err != nil {
			http.Error(w, "invalid time", http.StatusBadRequest)
			return
		}

		transactions, err := h.adminModule.Actions.GetTransactionHistory.Do(ctx, studentID, from, to)
		if err != nil {
			http.Error(w, "failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := h.prepareResponse(transactions)

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) prepareResponse(transactions []dto.TransactionHistory) Response {
	return Response{
		Transactions: lo.Map(transactions, func(t dto.TransactionHistory, _ int) Transaction {
			loc, _ := time.LoadLocation("Europe/Moscow")
			return Transaction{
				ID:          t.ID.String(),
				CreatedAt:   t.CreatedAt.In(loc).String(),
				Amount:      t.Amount.String(),
				IsConfirmed: t.IsConfirmed,
			}
		}),
		TransactionsCount: int64(len(transactions)),
	}
}
