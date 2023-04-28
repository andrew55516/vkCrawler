package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"time"
	db "vkCrawler/db/sqlc"
)

var sheets = []string{
	"LikesOnlyUsers",
	"NewLikesOnlyUsers",
	"BecameCommenter",
	"AmountOfLikesOnlyUsers",
	"AmountOfNewLikesOnlyUsers",
	"AmountOfBecameCommenter"}

func LikesOnlyUsersToXLSX() error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for _, sheet := range sheets {
		_, err := f.NewSheet(sheet)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	lastUsers := make(map[string]bool, 0)

	location, _ := time.LoadLocation("Europe/Moscow")
	layout := "2006-01-02"
	startTime, _ := time.ParseInLocation(layout, "2020-03-22", location)

	for i := 1; i < 37; i++ {
		toTime := startTime.AddDate(0, i, 0)
		users, err := db.GetLikesOnlyUsers(toTime)
		if err != nil {
			return err
		}

		newUsers := make([]string, 0)
		becameCommenter := make([]string, 0)

		for _, user := range users {
			if _, ok := lastUsers[user]; ok {
				lastUsers[user] = true
			} else {
				newUsers = append(newUsers, user)
			}
		}

		for k, v := range lastUsers {
			if !v {
				becameCommenter = append(becameCommenter, k)
			}
		}

		f.SetCellValue("LikesOnlyUsers", cols[i-1]+"1", toTime.String())
		for j, user := range users {
			f.SetCellValue("LikesOnlyUsers", cols[i-1]+fmt.Sprintf("%d", j+2), user)
		}

		f.SetCellValue("NewLikesOnlyUsers", cols[i-1]+"1", toTime.String())
		for j, user := range newUsers {
			f.SetCellValue("NewLikesOnlyUsers", cols[i-1]+fmt.Sprintf("%d", j+2), user)
		}

		f.SetCellValue("BecameCommenter", cols[i-1]+"1", toTime.String())
		for j, user := range becameCommenter {
			f.SetCellValue("BecameCommenter", cols[i-1]+fmt.Sprintf("%d", j+2), user)
		}

		f.SetCellValue("AmountOfLikesOnlyUsers", cols[i-1]+"1", toTime.String())
		f.SetCellValue("AmountOfLikesOnlyUsers", cols[i-1]+"2", len(users))

		f.SetCellValue("AmountOfNewLikesOnlyUsers", cols[i-1]+"1", toTime.String())
		f.SetCellValue("AmountOfNewLikesOnlyUsers", cols[i-1]+"2", len(newUsers))

		f.SetCellValue("AmountOfBecameCommenter", cols[i-1]+"1", toTime.String())
		f.SetCellValue("AmountOfBecameCommenter", cols[i-1]+"2", len(becameCommenter))

		lastUsers = make(map[string]bool, 0)
		for _, user := range users {
			lastUsers[user] = false
		}

	}

	if err := f.SaveAs("LikesOnlyUsers.xlsx"); err != nil {
		return err
	}

	return nil
}

var cols = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ"}
