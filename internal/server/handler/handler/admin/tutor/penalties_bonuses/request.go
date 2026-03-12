package penalties_bonuses

import (
	"strings"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/pkg/errors"
)

type Request struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Type    string `json:"type"`
	Amount  int64  `json:"amount"`
	Comment string `json:"comment"`
}

func (r Request) IsCreate() bool {
	return strings.TrimSpace(r.Type) != "" || r.Amount != 0 || strings.TrimSpace(r.Comment) != ""
}

func (r Request) ToTime() (time.Time, time.Time, error) {
	return convert.StringsIntervalToTime(r.From, r.To)
}

func (r Request) ToActualType() (dto.AccrualActualType, error) {
	switch strings.ToLower(strings.TrimSpace(r.Type)) {
	case "penalty", "штраф":
		return dto.AccrualActualTypePenalty, nil
	case "bonus", "премия":
		return dto.AccrualActualTypeBonus, nil
	default:
		return 0, errors.New("unsupported type")
	}
}
