package get_lessons

import (
	"time"

	"github.com/it-chep/tutors.git/internal/pkg/convert"
)

type Request struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (r Request) ToTime() (time.Time, time.Time, error) {
	return convert.StringsIntervalToTime(r.From, r.To)
}
