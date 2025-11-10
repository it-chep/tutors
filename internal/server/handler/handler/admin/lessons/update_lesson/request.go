package update_lesson

import (
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dto"
	"time"
)

type Request struct {
	Date     string `json:"date"`
	Duration int64  `json:"duration"`
}

func (r Request) ToDto() (dto.UpdateLesson, error) {
	createdDate, err := time.Parse(time.DateTime, r.Date)
	if err != nil {
		return dto.UpdateLesson{}, err
	}

	return dto.UpdateLesson{
		Date:     createdDate,
		Duration: time.Duration(r.Duration) * time.Minute,
	}, nil
}
