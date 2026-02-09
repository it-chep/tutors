package student_by_id

import (
	"encoding/json"
	"github.com/it-chep/tutors.git/internal/pkg/payment_hash"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
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

		studentIDStr := chi.URLParam(r, "student_id")
		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		baseData, err := h.adminModule.Actions.GetStudentByID.Do(ctx, studentID)
		if err != nil {
			http.Error(w, "failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := h.prepareResponse(baseData)

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) prepareResponse(student dto.Student) Response {
	var (
		paymentURL string
		err        error
	)
	if len(student.PaymentUUID) != 0 {
		paymentURL, err = payment_hash.EncryptPaymentData(student.ID, student.PaymentUUID)
		if err != nil {
			paymentURL = ""
		}
	}

	return Response{
		Student: Student{
			ID:                  student.ID,
			FirstName:           student.FirstName,
			LastName:            student.LastName,
			MiddleName:          student.MiddleName,
			Phone:               student.Phone,
			Tg:                  student.Tg,
			CostPerHour:         student.CostPerHour,
			SubjectName:         student.SubjectName,
			TutorName:           student.TutorName,
			TutorID:             student.TutorID,
			ParentFullName:      student.ParentFullName,
			ParentPhone:         student.ParentPhone,
			ParentTg:            student.ParentTg,
			Balance:             student.Balance.String(),
			IsOnlyTrialFinished: student.IsOnlyTrialFinished,
			IsBalanceNegative:   student.IsBalanceNegative,
			IsNewbie:            student.IsNewbie,
			ParentTgID:          student.ParentTgID,
			TgAdminUsername:     student.TgAdminUsername,
			IsArchived:          student.IsArchived,

			PaymentName: student.Payment.String(),
			PaymentID:   student.Payment.ID,
			PaymentURL:  paymentURL,
		},
	}
}
