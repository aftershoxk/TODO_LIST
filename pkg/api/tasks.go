package api

import (
	"net/http"
	"todo_list/pkg/db"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(tasksLimit, search)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
