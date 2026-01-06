package api

import (
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
		return
	case http.MethodGet:
		getTaskHandler(w, r)
		return
	case http.MethodPut:
		updateTaskHandler(w, r)
		return
	case http.MethodDelete:
		deleteTaskHandler(w, r)
		return
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method error")
	}
}
