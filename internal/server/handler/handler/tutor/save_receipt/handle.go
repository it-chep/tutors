package save_receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/payout"
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

		file, header, contentType, err := extractMultipartFile(r, "receipt", "file")
		if err != nil {
			http.Error(w, "failed to read file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		payoutID, err := h.adminModule.Actions.Payouts.SaveReceipt(ctx, payout.ReceiptUploadRequest{
			FileName:    header.Filename,
			ContentType: contentType,
			Body:        file,
		})
		if err != nil {
			http.Error(w, "failed to save receipt: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": payoutID.String(),
		})
	}
}

type multipartHeader struct {
	Filename string
}

func extractMultipartFile(r *http.Request, keys ...string) (multipartFile io.ReadCloser, header *multipartHeader, contentType string, err error) {
	if err = r.ParseMultipartForm(32 << 20); err != nil {
		return nil, nil, "", err
	}

	for _, key := range keys {
		file, fh, openErr := r.FormFile(key)
		if openErr != nil {
			continue
		}

		return file, &multipartHeader{Filename: fh.Filename}, fh.Header.Get("Content-Type"), nil
	}

	return nil, nil, "", fmt.Errorf("file field not found")
}
