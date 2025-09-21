package dao

import (
	"database/sql"
	"strings"

	"github.com/it-chep/tutors.git/internal/module/admin/dto"
)

type TutorDAO struct {
	TutorID     sql.NullInt64 `db:"id"`
	CostPerHour string        `db:"cost_per_hour" json:"cost_per_hour"`
	SubjectID   int64         `db:"subject_id"`
	AdminID     int64         `db:"admin_id"`
	FullName    string        `db:"full_name"`
	Tg          string        `db:"tg"`
	Phone       string        `db:"phone"`
}

type TutorWithSubjectName struct {
	TutorDAO
	SubjectName sql.NullString `db:"subject_name"`
}

func (t TutorDAO) ToDomain() dto.Tutor {
	return dto.Tutor{
		ID:          t.TutorID.Int64,
		FullName:    t.FullName,
		Phone:       t.Phone,
		Tg:          TgLink(t.Tg),
		CostPerHour: t.CostPerHour,
		SubjectID:   t.SubjectID,
		AdminID:     t.AdminID,
	}
}

type TutorsDao []TutorDAO

func (ts TutorsDao) ToDomain() []dto.Tutor {
	domain := make([]dto.Tutor, 0, len(ts))
	for _, t := range ts {
		domain = append(domain, t.ToDomain())
	}
	return domain
}

type TutorLessonsCountDao struct {
	LessonsCount int64 `db:"lessons_count" json:"lessons_count"`
	TrialCount   int64 `db:"trial_count" json:"trial_lessons"`
	BaseCount    int64 `db:"base_count" json:"base_lessons"`
}

func (l TutorLessonsCountDao) ToDomain() dto.TutorLessons {
	return dto.TutorLessons{
		TrialCount:   l.TrialCount,
		BaseCount:    l.BaseCount,
		LessonsCount: l.LessonsCount,
	}
}

func TgLink(tg string) string {
	if tg == "" {
		return ""
	}

	// Удаляем @ в начале, если есть
	tgURL := strings.TrimPrefix(tg, "@")

	// Если URL не содержит http/https, формируем полный URL
	if !strings.HasPrefix(tgURL, "http") {
		tgURL = "https://t.me/" + tgURL
	}

	return tgURL
}
