package register

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/dto"
	register_dal "github.com/it-chep/tutors.git/internal/module/admin/action/auth/register/dal"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/it-chep/tutors.git/pkg/cache"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/it-chep/tutors.git/pkg/token"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Action struct {
	jwt   config.JwtConfig
	codes *cache.Cache[string, register_dto.CodeRegister]
	repo  *register_dal.Repository
	smtp  *smtp.ClientSmtp
}

func New(pool *pgxpool.Pool, smtp *smtp.ClientSmtp, jwt config.JwtConfig) *Action {
	return &Action{
		jwt:   jwt,
		codes: cache.NewCache[string, register_dto.CodeRegister](1000, time.Minute),
		repo:  register_dal.NewRepository(pool),
		smtp:  smtp,
	}
}

func (a *Action) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req register_dto.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := req.Validate(); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		exists, err := a.repo.IsEmailExists(r.Context(), req.Email)
		if err != nil {
			logger.Error(r.Context(), "произошла ошибка проверки email", err)
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Не смогли найти такого email", http.StatusBadRequest)
			return
		}

		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			return
		}

		code := smtp.GenerateCode()
		a.codes.Put(req.Email, register_dto.CodeRegister{Password: string(passHash), Code: code})
		err = a.smtp.SendEmail(smtp.EmailParams{
			Body: fmt.Sprintf("Ваш код %d", code), Destination: req.Email,
			Subject: "Регистрация в системе 100rep.ru",
		})
		if err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (a *Action) VerifyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req register_dto.VerifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		code, ok := a.codes.Get(req.Email)
		defer a.codes.Remove(req.Email)
		if !ok || code.Code != req.Code {
			http.Error(w, "invalid code", http.StatusBadRequest)
			return
		}

		if err := a.repo.SavePass(r.Context(), req.Email, code.Password); err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			return
		}

		tokens, err := token.GenerateTokens(req.Email, a.jwt.JwtSecret, a.jwt.RefreshSecret)
		if err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{
			Name:     register_dto.RefreshCookie,
			Value:    tokens.Refresh(),
			Expires:  time.Now().UTC().Add(14 * 24 * time.Hour),
			HttpOnly: true,
		})
		_ = json.NewEncoder(w).Encode(tokens)
	}
}
