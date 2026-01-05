package api

import (
	"net/http"
	"time"
	"todo_list/pkg/db"
	rep "todo_list/pkg/repeat"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "error: method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, 400, "error: invalid data")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, 400, err.Error())
			return
		}
	} else {
		now, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			writeError(w, 400, err.Error())
			return
		}
		today := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			0, 0, 0, 0,
			now.Location(),
		)
		next, err := rep.NextDate(today, task.Date, task.Repeat)
		if err != nil {
			writeError(w, 400, err.Error())
			return
		}
		err = db.UpdateDate(id, next)
		if err != nil {
			writeError(w, 400, err.Error())
			return
		}
	}
	writeJSON(w, 200, struct{}{})
}
