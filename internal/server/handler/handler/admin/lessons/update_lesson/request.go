package update_lesson

import (
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dto"
	"strconv"
	"time"
)

type Request struct {
	Date     string `json:"date"`
	Duration string `json:"duration"`
}

func (r Request) ToDto() (dto.UpdateLesson, error) {
	createdDate, err := time.Parse(time.DateTime, r.Date)
	if err != nil {
		return dto.UpdateLesson{}, err
	}

	dura, _ := strconv.Atoi(r.Duration)
	return dto.UpdateLesson{
		Date:     createdDate,
		Duration: time.Duration(dura) * time.Minute,
	}, nil
}
