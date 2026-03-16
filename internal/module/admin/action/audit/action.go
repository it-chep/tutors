package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type routeKey struct {
	Method string
	Path   string
}

type meta struct {
	Description  string
	Action       string
	EntityName   string
	EntityParam  string
	CaptureState bool
}

type Entry struct {
	UserID      int64
	Description string
	Body        *string
	Action      string
	EntityName  string
	EntityID    *int64
}

type creator interface {
	Create(ctx context.Context, entry Entry) error
	Snapshot(ctx context.Context, entityName string, entityID int64) (map[string]any, error)
}

type Action struct {
	repo creator
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		repo: NewRepository(pool),
	}
}

func (a *Action) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			meta, ok := auditRoutes[routeKey{
				Method: r.Method,
				Path:   normalizePath(r.URL.Path),
			}]
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			body, err := readBody(r)
			if err != nil {
				logger.Error(r.Context(), "failed to read request body for audit", err)
				next.ServeHTTP(w, r)
				return
			}

			entityID := extractEntityID(r, meta.EntityParam)
			var before map[string]any
			if meta.CaptureState && entityID != nil {
				before, err = a.repo.Snapshot(r.Context(), meta.EntityName, *entityID)
				if err != nil {
					logger.Error(r.Context(), "failed to fetch audit snapshot before request", err)
				}
			}

			rec := &statusRecorder{ResponseWriter: w}
			next.ServeHTTP(rec, r)

			if rec.Status() >= http.StatusBadRequest {
				return
			}

			userID := userCtx.UserIDFromContext(r.Context())
			if userID == 0 {
				return
			}

			entry := Entry{
				UserID:      userID,
				Description: meta.Description,
				Action:      meta.Action,
				EntityName:  meta.EntityName,
				EntityID:    entityID,
			}

			var after map[string]any
			if meta.CaptureState && entry.EntityID != nil {
				after, err = a.repo.Snapshot(r.Context(), entry.EntityName, *entry.EntityID)
				if err != nil {
					logger.Error(r.Context(), "failed to fetch audit snapshot after request", err)
				}
			}

			entry.Body = buildBody(body, before, after)

			if err = a.repo.Create(r.Context(), entry); err != nil {
				logger.Error(r.Context(), "failed to write audit log", err)
			}
		})
	}
}

func (a *Action) Create(ctx context.Context, entry Entry) error {
	return a.repo.Create(ctx, entry)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *statusRecorder) Write(body []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(body)
}

func (r *statusRecorder) Status() int {
	if r.status == 0 {
		return http.StatusOK
	}
	return r.status
}

func readBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

func buildBody(body []byte, before, after map[string]any) *string {
	payload := map[string]any{}

	if request := requestBodyValue(body); request != nil {
		payload["request"] = request
	}
	if before != nil {
		payload["before"] = before
	}
	if after != nil {
		payload["after"] = after
	}
	if changes := buildChanges(before, after); len(changes) > 0 {
		payload["changes"] = changes
	}
	if len(payload) == 0 {
		return nil
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil
	}

	s := string(raw)
	return &s
}

func requestBodyValue(body []byte) any {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return nil
	}

	var value any
	if err := json.Unmarshal(trimmed, &value); err == nil {
		return value
	}

	return map[string]string{
		"raw_body": string(trimmed),
	}
}

func buildChanges(before, after map[string]any) map[string]map[string]any {
	if before == nil && after == nil {
		return nil
	}

	keys := map[string]struct{}{}
	for key := range before {
		keys[key] = struct{}{}
	}
	for key := range after {
		keys[key] = struct{}{}
	}

	changes := make(map[string]map[string]any)
	for key := range keys {
		var beforeValue any
		var afterValue any
		if before != nil {
			beforeValue = before[key]
		}
		if after != nil {
			afterValue = after[key]
		}

		if reflect.DeepEqual(beforeValue, afterValue) {
			continue
		}

		changes[key] = map[string]any{
			"before": beforeValue,
			"after":  afterValue,
		}
	}

	return changes
}

func extractEntityID(r *http.Request, param string) *int64 {
	if param == "" {
		return nil
	}

	raw := chi.URLParam(r, param)
	if raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err == nil {
			return &id
		}
	}

	match := regexp.MustCompile(`/(\d+)`).FindStringSubmatch(r.URL.Path)
	if len(match) != 2 {
		return nil
	}

	id, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return nil
	}

	return &id
}

func normalizePath(path string) string {
	re := regexp.MustCompile(`/\d+`)
	normalized := re.ReplaceAllString(path, `/{id}`)

	if strings.HasSuffix(normalized, "/{id") && !strings.HasSuffix(normalized, "/{id}") {
		normalized += "}"
	}

	return normalized
}

var auditRoutes = map[routeKey]meta{
	{Method: http.MethodPost, Path: "/admin/admins"}: {
		Description: "Создание администратора",
		Action:      "Создание администратора",
		EntityName:  "admin",
	},
	{Method: http.MethodDelete, Path: "/admin/admins/{id}"}: {
		Description:  "Удаление администратора",
		Action:       "Удаление администратора",
		EntityName:   "admin",
		EntityParam:  "admin_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/assistant"}: {
		Description: "Создание ассистента",
		Action:      "Создание ассистента",
		EntityName:  "assistant",
	},
	{Method: http.MethodDelete, Path: "/admin/assistant/{id}"}: {
		Description:  "Удаление ассистента",
		Action:       "Удаление ассистента",
		EntityName:   "assistant",
		EntityParam:  "assistant_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/assistant/{id}/add_available_tg"}: {
		Description:  "Добавление доступных Telegram аккаунтов ассистенту",
		Action:       "Добавление TG ассистенту",
		EntityName:   "assistant",
		EntityParam:  "assistant_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/assistant/{id}/permissions"}: {
		Description:  "Обновление прав ассистента",
		Action:       "Обновление прав ассистента",
		EntityName:   "assistant",
		EntityParam:  "assistant_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/assistant/{id}/delete_available_tg"}: {
		Description:  "Удаление доступных Telegram аккаунтов ассистента",
		Action:       "Удаление TG у ассистента",
		EntityName:   "assistant",
		EntityParam:  "assistant_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/assistant/{id}/penalties-bonuses"}: {
		Description: "Начисление штрафа или премии ассистенту",
		Action:      "Начисление ассистенту",
		EntityName:  "assistant",
		EntityParam: "assistant_id",
	},
	{Method: http.MethodPost, Path: "/admin/tutors"}: {
		Description: "Создание репетитора",
		Action:      "Создание репетитора",
		EntityName:  "tutor",
	},
	{Method: http.MethodDelete, Path: "/admin/tutors/{id}"}: {
		Description:  "Удаление репетитора",
		Action:       "Удаление репетитора",
		EntityName:   "tutor",
		EntityParam:  "tutor_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/tutors/trial_lesson"}: {
		Description: "Проведение пробного занятия",
		Action:      "Проведение пробного занятия",
		EntityName:  "lesson",
	},
	{Method: http.MethodPost, Path: "/admin/tutors/conduct_lesson"}: {
		Description: "Проведение занятия",
		Action:      "Проведение занятия",
		EntityName:  "lesson",
	},
	{Method: http.MethodPost, Path: "/admin/tutors/{id}/archive"}: {
		Description:  "Архивация репетитора",
		Action:       "Архивация репетитора",
		EntityName:   "tutor",
		EntityParam:  "tutor_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/tutors/{id}/unarchive"}: {
		Description:  "Разархивация репетитора",
		Action:       "Разархивация репетитора",
		EntityName:   "tutor",
		EntityParam:  "tutor_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/tutors/{id}/update"}: {
		Description:  "Обновление репетитора",
		Action:       "Обновление репетитора",
		EntityName:   "tutor",
		EntityParam:  "tutor_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/tutors/{id}/contract"}: {
		Description: "Загрузка договора репетитора",
		Action:      "Загрузка договора",
		EntityName:  "tutor",
		EntityParam: "tutor_id",
	},
	{Method: http.MethodDelete, Path: "/admin/tutors/{id}/contract"}: {
		Description: "Удаление договора репетитора",
		Action:      "Удаление договора",
		EntityName:  "tutor",
		EntityParam: "tutor_id",
	},
	{Method: http.MethodPost, Path: "/admin/tutors/{id}/penalties-bonuses"}: {
		Description: "Начисление штрафа или премии репетитору",
		Action:      "Начисление репетитору",
		EntityName:  "tutor",
		EntityParam: "tutor_id",
	},
	{Method: http.MethodPost, Path: "/admin/tutors/{id}/payouts"}: {
		Description: "Создание выплаты репетитору",
		Action:      "Создание выплаты",
		EntityName:  "tutor",
		EntityParam: "tutor_id",
	},
	{Method: http.MethodPost, Path: "/admin/students/push_all_students"}: {
		Description: "Массовая отправка уведомлений студентам",
		Action:      "Массовая отправка уведомлений",
		EntityName:  "notification",
	},
	{Method: http.MethodPost, Path: "/admin/students"}: {
		Description: "Создание студента",
		Action:      "Создание студента",
		EntityName:  "student",
	},
	{Method: http.MethodPost, Path: "/admin/students/move"}: {
		Description: "Перемещение студентов",
		Action:      "Перемещение студентов",
		EntityName:  "student",
	},
	{Method: http.MethodPost, Path: "/admin/students/change_all_payment"}: {
		Description: "Массовая смена платёжной системы студентов",
		Action:      "Массовая смена платёжки",
		EntityName:  "student",
	},
	{Method: http.MethodDelete, Path: "/admin/students/{id}"}: {
		Description:  "Удаление студента",
		Action:       "Удаление студента",
		EntityName:   "student",
		EntityParam:  "student_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/students/{id}"}: {
		Description:  "Обновление студента",
		Action:       "Обновление студента",
		EntityName:   "student",
		EntityParam:  "student_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/students/{id}/wallet"}: {
		Description:  "Обновление кошелька студента",
		Action:       "Обновление кошелька",
		EntityName:   "wallet",
		EntityParam:  "student_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/students/{id}/change_payment"}: {
		Description:  "Смена платёжной системы студента",
		Action:       "Смена платёжной системы",
		EntityName:   "student",
		EntityParam:  "student_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/students/{id}/transactions/manual"}: {
		Description:  "Добавление ручной транзакции студенту",
		Action:       "Добавление ручной транзакции",
		EntityName:   "student",
		EntityParam:  "student_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/students/{id}/notifications/push"}: {
		Description: "Отправка уведомления студенту",
		Action:      "Отправка уведомления",
		EntityName:  "student",
		EntityParam: "student_id",
	},
	{Method: http.MethodPost, Path: "/admin/students/{id}/archive"}: {
		Description:  "Архивация студента",
		Action:       "Архивация студента",
		EntityName:   "student",
		EntityParam:  "student_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/students/{id}/unarchive"}: {
		Description:  "Разархивация студента",
		Action:       "Разархивация студента",
		EntityName:   "student",
		EntityParam:  "student_id",
		CaptureState: true,
	},
	{Method: http.MethodDelete, Path: "/admin/lessons/{id}"}: {
		Description:  "Удаление урока",
		Action:       "Удаление урока",
		EntityName:   "lesson",
		EntityParam:  "lesson_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/lessons/{id}"}: {
		Description:  "Обновление урока",
		Action:       "Обновление урока",
		EntityName:   "lesson",
		EntityParam:  "lesson_id",
		CaptureState: true,
	},
	{Method: http.MethodPost, Path: "/admin/students/{id}/comments"}: {
		Description: "Создание комментария к студенту",
		Action:      "Создание комментария",
		EntityName:  "student",
		EntityParam: "student_id",
	},
	{Method: http.MethodDelete, Path: "/admin/students/{id}/comments/{id}"}: {
		Description: "Удаление комментария студента",
		Action:      "Удаление комментария",
		EntityName:  "student",
		EntityParam: "student_id",
	},
}
