package logout

import (
	register_dto "github.com/it-chep/tutors.git/internal/module/admin/action/auth/dto"
	"net/http"
	"time"
)

type Action struct{}

func New() *Action {
	return &Action{}
}

func (a *Action) DeleteCookieHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     register_dto.RefreshCookie,
			Value:    "",
			Expires:  time.Now().UTC().Add(-1 * time.Hour),
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})

		// Дополнительно можно отправить ответ
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cookie deleted"))
	}
}
