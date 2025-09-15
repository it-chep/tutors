package convert

import (
	"github.com/pkg/errors"
	"time"
)

func StringsIntervalToTime(from, to string) (time.Time, time.Time, error) {
	fromTime, err := time.Parse(time.DateOnly, from)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Неправильно указан формат даты 'ОТ'")
	}

	toTime, err := time.Parse(time.DateOnly, to)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Неправильно указан формат даты 'ДО'")
	}
	// Используем текущую временную зону
	loc := time.Now().Location()

	// Устанавливаем время для fromTime: 00:00:00 в текущей локации
	fromTime = time.Date(
		fromTime.Year(),
		fromTime.Month(),
		fromTime.Day(),
		0, 0, 0, 0,
		loc,
	).Add(24 * time.Hour)

	// Устанавливаем время для toTime: 23:59:59 в текущей локации
	toTime = time.Date(
		toTime.Year(),
		toTime.Month(),
		toTime.Day(),
		23, 59, 59, 0,
		loc,
	).Add(24 * time.Hour)

	return fromTime, toTime, nil
}
