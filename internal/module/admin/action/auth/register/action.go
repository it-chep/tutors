package register

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/it-chep/tutors.git/internal/config"
	register_dto "github.com/it-chep/tutors.git/internal/module/admin/action/auth/dto"
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
		codes: cache.NewCache[string, register_dto.CodeRegister](1000, 15*time.Minute),
		repo:  register_dal.NewRepository(pool),
		smtp:  smtp,
	}
}

func (a *Action) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req register_dto.RegisterRequest
		ctx := r.Context()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			logger.Error(r.Context(), "rega Ошибка при декоде", err)
			return
		}

		if err := req.Validate(); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			logger.Error(r.Context(), "Rega Ошибка при валидации", err)
			return
		}

		exists, err := a.repo.IsEmailExists(ctx, req.Email)
		if err != nil {
			logger.Error(ctx, "произошла ошибка проверки email", err)
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Не смогли найти такого email", http.StatusBadRequest)
			logger.Error(r.Context(), "Rega нет юзера", err)
			return
		}

		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			logger.Error(r.Context(), "Rega пароль ошибка", err)
			return
		}

		code := smtp.GenerateCode()
		//fmt.Println("code: ", code)
		a.codes.Put(req.Email, register_dto.CodeRegister{Password: string(passHash), Code: code})
		if os.Getenv("DEBUG") != "True" {
			err = a.smtp.SendEmail(smtp.EmailParams{
				Body: fmt.Sprintf("Ваш код %s", code), Destination: req.Email,
				Subject: "Регистрация в системе 100rep.ru",
			})
			if err != nil {
				http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
				logger.Error(r.Context(), "Ошибка при отправке кода", err)
				return
			}
		}

		err = a.repo.SaveCode(ctx, req.Email, code)
		if err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			logger.Error(r.Context(), "Ошибка при сохранении кода", err)
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
			logger.Error(r.Context(), "Ошибка при сохранении декоде", err)
			return
		}

		code, ok := a.codes.Get(req.Email)
		defer a.codes.Remove(req.Email)
		if !ok || code.Code != req.Code {
			http.Error(w, "invalid code", http.StatusBadRequest)
			logger.Error(r.Context(), "Ошибка при коде", nil)
			return
		}

		if err := a.repo.SavePass(r.Context(), req.Email, code.Password); err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			logger.Error(r.Context(), "Ошибка при сохр пароля", err)
			return
		}

		tokens, err := token.GenerateTokens(req.Email, a.jwt.JwtSecret, a.jwt.RefreshSecret)
		if err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			logger.Error(r.Context(), "Ошибка при токене", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{
			Name:     register_dto.RefreshCookie,
			Value:    tokens.Refresh(),
			Expires:  time.Now().UTC().Add(60 * 24 * time.Hour),
			HttpOnly: true,
		})
		_ = json.NewEncoder(w).Encode(tokens)
	}
}
