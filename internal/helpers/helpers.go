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
	str = strings.Replace(str, " ", " ", -1)
	parts := strings.Split(str, " ")
	switch parts[0] {
	case "сегодня":
		timeStr = fmt.Sprintf("2023-03-21 %s", parts[2])
	case "вчера":
		timeStr = fmt.Sprintf("2023-03-20 %s", parts[2])
	default:
		if len(parts[0]) == 1 {
			parts[0] = "0" + parts[0]
		}
		switch len(parts) {
		case 5:
			timeStr = fmt.Sprintf("%s-%s-%s %s", parts[2], months[parts[1]], parts[0], parts[4])
		case 4:
			timeStr = fmt.Sprintf("2023-%s-%s %s", months[parts[1]], parts[0], parts[3])
		case 3:
			timeStr = fmt.Sprintf("%s-%s-%s 12:00", parts[2], months[parts[1]], parts[0])
		}

	}
	t, _ := time.ParseInLocation(layout, timeStr, location)
	return t
}
