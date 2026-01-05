package api

import (
	"encoding/json"
	"net/http"
	"todo_list/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if task.ID == "" {
		writeError(w, 404, "id is nil")
		return
	}
	if task.Title == "" {
		writeError(w, 400, "invalid data")
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeError(w, 400, "invalid data")
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, db.Task{})
}
