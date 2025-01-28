package tasks

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Функция расчета следующей даты задачи
func NextDate(now time.Time, date string, repeat string) (string, error) {

	tDate, err := time.Parse("20060102", date) // Преобразование из строки в дату
	if err != nil {
		return "", errors.New("ошибка формата даты")
	}

	switch repeat {

	case "": // Задача без повторов
		return "", nil

	case "y": // Повторять задачу ежегодно
		tDate = tDate.AddDate(1, 0, 0)
		for tDate.Before(now) {
			tDate = tDate.AddDate(1, 0, 0)
		}
		return tDate.Format("20060102"), nil

	default:
		parts := strings.Fields(repeat)
		if len(parts) == 0 {
			return "", errors.New("неверный формат повторения")
		}

		// Повтор задачи через d дней
		if strings.HasPrefix(repeat, "d ") && len(parts) == 2 {
			days, err := strconv.Atoi(parts[1])
			if err != nil || days < 1 || days > 400 {
				return "", fmt.Errorf("недопустимое количество дней: %s", repeat)
			}
			tDate = tDate.AddDate(0, 0, days)
			for tDate.Before(now) {
				tDate = tDate.AddDate(0, 0, days)
			}
			return tDate.Format("20060102"), nil
		}

		// w Повтор задачи назначается в указанные дни недели
		if strings.HasPrefix(repeat, "w ") && len(parts) == 2 {

			if now.After(tDate) {
				tDate = now
			}

			var nextDate string
			wDays := strings.Split(parts[1], ",")
			for i := range wDays {
				day, err := strconv.Atoi(wDays[i])
				if err != nil || !(day > 0 && day < 8) {
					return "", fmt.Errorf("недопустимый день недели: %s", repeat)
				}

				dateDay := tDate
				dateDay = dateDay.AddDate(0, 0, 1)
				dayofWeek := time.Weekday(day % 7)
				dateWeekday := dateDay.Weekday()

				for dateWeekday != dayofWeek {
					dateDay = dateDay.AddDate(0, 0, 1)
					dateWeekday = dateDay.Weekday()
				}

				if nextDate == "" {
					nextDate = dateDay.Format("20060102")
				} else if nextDate > dateDay.Format("20060102") {
					nextDate = dateDay.Format("20060102")
				}
			}
			return nextDate, nil
		}

		// m Повтор задачи назначается в указанные дни месяцев
		if strings.HasPrefix(repeat, "m ") && len(parts) > 1 {
			if now.After(tDate) {
				tDate = now
			}

			var monthsStr []string
			var setMonth bool
			mDays := strings.Split(parts[1], ",")
			if len(parts) == 3 {
				monthsStr = strings.Split(parts[2], ",")
				setMonth = true
			}

			var days []int
			for _, val := range mDays {
				day, err := strconv.Atoi(val)
				if err != nil || !(day > -3 && day < 32) {
					return "", fmt.Errorf("недопустимый день недели: %s", repeat)
				}
				days = append(days, day)
			}

			dateDay := tDate
			var nextDate string
			if !setMonth {
				for _, day := range days {
					dateMonth := tDate.Month()

					if day == -1 || day == -2 { //Последний или предпоследний день месяца
						dateDay = time.Date(tDate.Year(), tDate.Month()+1, 1, 0, 0, 0, 0, tDate.Location()).AddDate(0, 0, day)
					} else if day == 31 {
						if dateMonth == time.February || dateMonth == time.April || dateMonth == time.June || dateMonth == time.September || dateMonth == time.November {
							dateMonth = dateMonth + 1
						}

						dateDay = time.Date(tDate.Year(), dateMonth, day, 0, 0, 0, 0, tDate.Location())

					} else if day > 28 && dateMonth == time.February {
						dateDay = time.Date(tDate.Year(), dateMonth+1, 1, 0, 0, 0, 0, tDate.Location()).AddDate(0, 0, -1)
					} else {
						dateDay = time.Date(tDate.Year(), tDate.Month(), day, 0, 0, 0, 0, tDate.Location())
					}

					for !now.Before(dateDay) {
						dateDay = dateDay.AddDate(0, 1, 0)
					}

					if nextDate == "" {
						nextDate = dateDay.Format("20060102")
					} else if nextDate > dateDay.Format("20060102") {
						nextDate = dateDay.Format("20060102")
					}

				}
			} else {
				for _, val := range monthsStr {
					month, _ := strconv.Atoi(val)
					dateMonth := time.Month(month)

					for _, day := range days {
						if day == -1 || day == -2 { //Последний или предпоследний день месяца
							dateDay = time.Date(tDate.Year(), dateMonth+1, 1, 0, 0, 0, 0, tDate.Location()).AddDate(0, 0, day)
						} else if day == 31 {
							if dateMonth == time.February || dateMonth == time.April || dateMonth == time.June || dateMonth == time.September || dateMonth == time.November {
								dateMonth = dateMonth + 1
							}

							dateDay = time.Date(tDate.Year(), dateMonth, day, 0, 0, 0, 0, tDate.Location())

						} else if day > 28 && dateMonth == time.February {
							dateDay = time.Date(tDate.Year(), dateMonth+1, 1, 0, 0, 0, 0, tDate.Location()).AddDate(0, 0, -1)
						} else {
							dateDay = time.Date(tDate.Year(), dateMonth, day, 0, 0, 0, 0, tDate.Location())
						}

						for dateDay.Before(now) {
							dateDay = dateDay.AddDate(1, 0, 0)
						}

						if nextDate == "" {
							nextDate = dateDay.Format("20060102")

						} else if nextDate > dateDay.Format("20060102") {
							nextDate = dateDay.Format("20060102")
						}

					}

				}

			}
			return nextDate, nil

		}

	}
	return "", fmt.Errorf("неподдерживаемый формат повторения: %s", repeat)
}
