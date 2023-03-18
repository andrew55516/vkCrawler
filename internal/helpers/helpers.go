package helpers

import (
	"fmt"
	"strings"
	"time"
)

func StrToTime(str string) time.Time {
	location, _ := time.LoadLocation("Europe/Moscow")
	layout := "2006-01-02 15:04"
	months := map[string]string{
		"янв": "01",
		"фев": "02",
		"мар": "03",
		"апр": "04",
		"мая": "05",
		"июн": "06",
		"июл": "07",
		"авг": "08",
		"сен": "09",
		"окт": "10",
		"ноя": "11",
		"дек": "12",
	}

	var timeStr string
	parts := strings.Split(str, " ")
	switch parts[0] {
	case "сегодня":
		timeStr = fmt.Sprintf("2023-03-13 %s", parts[2])
	case "вчера":
		timeStr = fmt.Sprintf("2023-03-12 %s", parts[2])
	default:
		if len(parts[0]) == 1 {
			parts[0] = "0" + parts[0]
		}
		timeStr = fmt.Sprintf("%s-%s-%s %s", parts[2], months[parts[1]], parts[0], parts[4])

	}
	t, _ := time.ParseInLocation(layout, timeStr, location)
	return t
}
