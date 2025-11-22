package convert

import (
	"time"

	"github.com/pkg/errors"
)

func StringsIntervalToTime(from, to string) (time.Time, time.Time, error) {
	fromTime, err := time.Parse(time.DateTime, from)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Неправильно указан формат даты 'ОТ'")
	}

	toTime, err := time.Parse(time.DateTime, to)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Неправильно указан формат даты 'ДО'")
	}

	newToTime := toTime.Add(24 * time.Hour)

	return fromTime.UTC(), newToTime.UTC(), nil
}
