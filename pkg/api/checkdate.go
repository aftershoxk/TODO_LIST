package api

import (
	"time"
	"todo_list/pkg/db"
	n "todo_list/pkg/repeat"
)

func checkDate(task *db.Task) error {
	today := time.Now()

	todayDate := time.Date(
		today.Year(),
		today.Month(),
		today.Day(),
		0, 0, 0, 0,
		today.Location(),
	)

	if task.Date == "" {
		task.Date = todayDate.Format(DateFormat)
		return nil
	}

	parsed, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		if parsed.Before(todayDate) {
			task.Date = todayDate.Format(DateFormat)
		}
		return nil
	}

	if parsed.After(todayDate) || parsed.Equal(todayDate) {
		return nil
	}

	next, err := n.NextDate(todayDate, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	task.Date = next
	return nil
}
