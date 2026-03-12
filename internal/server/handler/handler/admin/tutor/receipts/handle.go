package receipts

import (
	"archive/zip"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/payout"
	"github.com/samber/lo"
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

		tutorID, err := strconv.ParseInt(chi.URLParam(r, "tutor_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid tutor ID", http.StatusBadRequest)
			return
		}

		var req Request
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		from, to, err := req.ToTime()
		if err != nil {
			http.Error(w, "invalid time", http.StatusBadRequest)
			return
		}

		items, err := h.adminModule.Actions.Payouts.ListReceipts(ctx, tutorID, from, to)
		if err != nil {
			http.Error(w, "failed to list receipts: "+err.Error(), http.StatusInternalServerError)
			return
		}
		mapped := mapItems(items)

		response := Response{
			Receipts: lo.Map(mapped, func(item payoutItem, _ int) Receipt {
				return Receipt{
					ID:                item.ID,
					Amount:            item.Amount,
					Comment:           item.Comment,
					CreatedAt:         item.CreatedAt,
					HasReceipt:        item.HasReceipt,
					ReceiptFileName:   item.ReceiptFileName,
					ReceiptUploadedAt: item.ReceiptUploadedAt,
				}
			}),
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) DownloadAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		from, to, err := req.ToTime()
		if err != nil {
			http.Error(w, "invalid time", http.StatusBadRequest)
			return
		}

		items, err := h.adminModule.Actions.Payouts.ListVisibleReceipts(ctx, from, to)
		if err != nil {
			http.Error(w, "failed to list receipts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\"receipts.zip\"")

		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

		for _, item := range items {
			if item.ReceiptFileName == "" {
				continue
			}

			downloaded, err := h.adminModule.Actions.Payouts.DownloadReceipt(ctx, item.ID)
			if err != nil {
				http.Error(w, "failed to download receipt: "+err.Error(), http.StatusInternalServerError)
				return
			}

			entry, err := zipWriter.Create(zipName(item.TutorName, item.ReceiptFileName))
			if err != nil {
				downloaded.Body.Close()
				http.Error(w, "failed to create zip entry: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err = io.Copy(entry, downloaded.Body); err != nil {
				downloaded.Body.Close()
				http.Error(w, "failed to write zip entry: "+err.Error(), http.StatusInternalServerError)
				return
			}
			downloaded.Body.Close()
		}
	}
}

type payoutItem struct {
	ID                string
	Amount            string
	Comment           string
	CreatedAt         string
	HasReceipt        bool
	ReceiptFileName   string
	ReceiptUploadedAt string
	TutorName         string
}

func zipName(parts ...string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_", " ", "_", ":", "_")
	return replacer.Replace(strings.Join(parts, "_"))
}

func mapItems(items []payout.Payout) []payoutItem {
	return lo.Map(items, func(item payout.Payout, _ int) payoutItem {
		uploadedAt := ""
		if item.ReceiptUploadedAt != nil {
			uploadedAt = item.ReceiptUploadedAt.Format(time.DateTime)
		}

		return payoutItem{
			ID:                item.ID.String(),
			Amount:            item.Amount.String(),
			Comment:           item.Comment,
			CreatedAt:         item.CreatedAt.Format(time.DateTime),
			HasReceipt:        item.ReceiptFileName != "",
			ReceiptFileName:   item.ReceiptFileName,
			ReceiptUploadedAt: uploadedAt,
			TutorName:         item.TutorName,
		}
	})
}
