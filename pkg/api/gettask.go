package api

import (
	"net/http"
	"todo_list/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, 400, "id is not defined")
		return
	}

	result, err := db.GetTask(id)
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}

	writeJSON(w, 200, result)
}
