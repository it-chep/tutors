package update_lesson

import (
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dto"
	"time"
)

type Request struct {
	Date     string `json:"date"`
	Duration int    `json:"duration"`
}

func (r Request) ToDto() (dto.UpdateLesson, error) {
	date, err := time.Parse(time.DateOnly, r.Date)
	if err != nil {
		return dto.UpdateLesson{}, err
	}

	return dto.UpdateLesson{
		Date:     date,
		Duration: time.Duration(r.Duration) * time.Minute,
	}, nil
}
