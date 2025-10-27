package login

import (
	"encoding/json"
	"fmt"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"net/http"
	"time"

	"github.com/it-chep/tutors.git/internal/config"
	register_dto "github.com/it-chep/tutors.git/internal/module/admin/action/auth/dto"
	login_dal "github.com/it-chep/tutors.git/internal/module/admin/action/auth/login/dal"
	"github.com/it-chep/tutors.git/pkg/cache"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/it-chep/tutors.git/pkg/token"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Action struct {
	jwt   config.JwtConfig
	codes *cache.Cache[string, string]
	repo  *login_dal.Repository
	smtp  *smtp.ClientSmtp
}

func New(pool *pgxpool.Pool, smtp *smtp.ClientSmtp, jwt config.JwtConfig) *Action {
	return &Action{
		jwt:   jwt,
		codes: cache.NewCache[string, string](1000, 15*time.Minute),
		repo:  login_dal.NewRepository(pool),
		smtp:  smtp,
	}
}

func (a *Action) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req register_dto.LoginRequest
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		user, err := a.repo.GetUser(r.Context(), req.Email)
		if err != nil {
			http.Error(w, "Не нашли такого пользователя", http.StatusBadRequest)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "неверный email или пароль", http.StatusBadRequest)
			return
		}

		code := smtp.GenerateCode()
		a.codes.Put(req.Email, code)
		err = a.smtp.SendEmail(smtp.EmailParams{
			Body: fmt.Sprintf("Ваш код %s", code), Destination: req.Email,
			Subject: "Авторизация в системе 100rep.ru",
		})
		if err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			logger.Error(r.Context(), "Ошибка при отправке кода", err)
			return
		}

		err = a.repo.SaveCode(r.Context(), req.Email, code)
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
			return
		}

		code, ok := a.codes.Get(req.Email)
		defer a.codes.Remove(req.Email)
		if !ok || code != req.Code {
			http.Error(w, "invalid code", http.StatusBadRequest)
			return
		}

		tokens, err := token.GenerateTokens(req.Email, a.jwt.JwtSecret, a.jwt.RefreshSecret)
		if err != nil {
			http.Error(w, "Пожалуйста, повторите попытку позже", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     register_dto.RefreshCookie,
			Value:    tokens.Refresh(),
			Expires:  time.Now().UTC().Add(60 * 24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(tokens)
	}
}
