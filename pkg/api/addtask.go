package api

import (
	"encoding/json"
	"net/http"
	"todo_list/pkg/db"
)

type response struct {
	ID int64 `json:"id"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, 400, "invalid data")
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, 400, "invalid data")
		return
	}
	resp := response{
		ID: id,
	}

	writeJSON(w, 200, resp)
}
