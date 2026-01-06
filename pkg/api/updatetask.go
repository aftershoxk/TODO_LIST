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
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if task.ID == "" {
		writeError(w, http.StatusNotFound, "id is nil")
		return
	}
	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "invalid data")
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid data")
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, db.Task{})
}
