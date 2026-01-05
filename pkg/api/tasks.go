package api

import (
	"net/http"
	"todo_list/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, 405, "method not allowed")
		return
	}

	limit := 50

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(limit, search)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, 200, TasksResp{Tasks: tasks})
}
