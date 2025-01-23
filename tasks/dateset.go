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
	tDate, err := time.Parse("20060102", date) //Преобразование из строки в дату
	if err != nil {
		return "", errors.New("ошибка формата даты")
	}

	switch repeat {
	case "": //Без повторов
		return "", errors.New("одноразовая задача")
	case "y": //Ежегодно
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

		//Повторить через d дней
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

		//Задача назначается в указанные дни недели w
		if strings.HasPrefix(repeat, "w ") && len(parts) == 2 {
			var multiDate string
			wDays := strings.Split(parts[1], ",")
			for i := range wDays {
				day, err := strconv.Atoi(wDays[i])
				if err != nil || !(day > 0 && day < 8) {
					return "", fmt.Errorf("недопустимый день недели: %s", repeat)
				}

				dayDate := tDate
				var found bool
				var i int
				found = false
				for !found {
					dayDate = dayDate.AddDate(0, 0, 1)
					dateWeekday := dayDate.Weekday()
					numWeekDay := int(dateWeekday) + 0
					i = i + 1
					if i == 8 {
						return "", fmt.Errorf("недопустимый формат дня недели: %s", repeat)
					}
					if numWeekDay == (day - 1) {
						dayDate = dayDate.AddDate(0, 0, 1)

						for dayDate.Before(now) || dayDate.Equal(now) {
							dayDate = dayDate.AddDate(0, 0, 7)
						}
						multiDate = multiDate + dayDate.Format("20060102") + " "
						found = true
					}

				}

			}
			return multiDate, nil

		}

		return "", fmt.Errorf("неподдерживаемый формат повторения: %s", repeat)
	}
}
