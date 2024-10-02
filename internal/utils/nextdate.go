package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if date == "" || repeat == "" {
		return "", errors.New("неверная дата или повторяющееся правило")
	}

	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("недопустимый формат даты: %w", err)
	}

	// проверка на ежедневный интервал
	if strings.HasPrefix(repeat, "d ") {
		daysStr := strings.TrimPrefix(repeat, "d ")
		days, err := strconv.Atoi(daysStr)
		fmt.Println(days)
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("недопустимый дневной интервал")
		}

		if isSameDate(startDate, now) && days == 1 {
			return startDate.Format("20060102"), nil
		}

		nextDate := startDate.AddDate(0, 0, days)
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}

		return nextDate.Format("20060102"), nil
	}

	// проверка на ежегодный интервал
	if repeat == "y" {
		nextDate := startDate.AddDate(1, 0, 0)

		// обработка 29 февраля
		if startDate.Month() == time.February && startDate.Day() == 29 {
			// проверка на високосный год следующего года
			if nextDate.Month() != time.February || nextDate.Day() != 29 {
				// переход на 1 марта, если след. год не високосный
				nextDate = time.Date(nextDate.Year(), time.March, 1, 0, 0, 0, 0, nextDate.Location())
			}
		}

		if nextDate.Before(now) {
			for !nextDate.After(now) {
				nextDate = nextDate.AddDate(1, 0, 0)
			}
		}
		return nextDate.Format("20060102"), nil
	}

	return "", errors.New("неподдерживаемый формат повтора")
}

func isSameDate(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}
