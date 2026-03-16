package contract

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	contractAction "github.com/it-chep/tutors.git/internal/module/admin/action/contract"
)

type Handler struct {
	adminModule *admin.Module
}

func NewHandler(adminModule *admin.Module) *Handler {
	return &Handler{
		adminModule: adminModule,
	}
}

func (h *Handler) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tutorID, err := strconv.ParseInt(chi.URLParam(r, "tutor_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid tutor ID", http.StatusBadRequest)
			return
		}

		file, header, contentType, err := extractMultipartFile(r, "file", "contract")
		if err != nil {
			http.Error(w, "failed to read file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		_, err = h.adminModule.Actions.Contracts.Upload(ctx, tutorID, contractAction.UploadRequest{
			FileName:    header.Filename,
			ContentType: contentType,
			Body:        file,
		})
		if err != nil {
			http.Error(w, "failed to upload contract: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) Download() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tutorID, err := strconv.ParseInt(chi.URLParam(r, "tutor_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid tutor ID", http.StatusBadRequest)
			return
		}

		contract, err := h.adminModule.Actions.Contracts.Get(ctx, tutorID)
		if err != nil {
			status := http.StatusInternalServerError
			if contractAction.IsNotFound(err) {
				status = http.StatusNotFound
			}
			http.Error(w, "failed to get contract: "+err.Error(), status)
			return
		}
		defer contract.Body.Close()

		contentType := contract.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", contract.FileName))
		if _, err = io.Copy(w, contract.Body); err != nil {
			http.Error(w, "failed to write contract: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tutorID, err := strconv.ParseInt(chi.URLParam(r, "tutor_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid tutor ID", http.StatusBadRequest)
			return
		}

		err = h.adminModule.Actions.Contracts.Delete(ctx, tutorID)
		if err != nil {
			status := http.StatusInternalServerError
			if contractAction.IsNotFound(err) {
				status = http.StatusNotFound
			}
			http.Error(w, "failed to delete contract: "+err.Error(), status)
			return
		}
	}
}

func (h *Handler) DownloadAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		contracts, err := h.adminModule.Actions.Contracts.ListVisible(ctx)
		if err != nil {
			http.Error(w, "failed to list contracts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\"contracts.zip\"")

		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

		for _, item := range contracts {
			downloaded, err := h.adminModule.Actions.Contracts.DownloadByKey(ctx, item.FileKey)
			if err != nil {
				http.Error(w, "failed to download contract: "+err.Error(), http.StatusInternalServerError)
				return
			}

			entry, err := zipWriter.Create(zipName(item.TutorName, item.FileName))
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

type multipartHeader struct {
	Filename string
}

func zipName(parts ...string) string {
	joined := strings.Join(parts, "_")
	replacer := strings.NewReplacer("/", "_", "\\", "_", " ", "_", ":", "_")
	return replacer.Replace(joined)
}

func (h *Handler) HandleMetadata() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tutorID, err := strconv.ParseInt(chi.URLParam(r, "tutor_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid tutor ID", http.StatusBadRequest)
			return
		}

		contract, err := h.adminModule.Actions.Contracts.Get(ctx, tutorID)
		if err != nil {
			status := http.StatusInternalServerError
			if contractAction.IsNotFound(err) {
				status = http.StatusNotFound
			}
			http.Error(w, "failed to get contract: "+err.Error(), status)
			return
		}
		contract.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"file_name":    contract.FileName,
			"content_type": contract.ContentType,
			"created_at":   contract.CreatedAt,
		})
	}
}
