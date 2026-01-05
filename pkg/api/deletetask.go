package api

import (
	"net/http"
	"todo_list/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, 400, "id is not defined")
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	writeJSON(w, 200, struct{}{})
}
