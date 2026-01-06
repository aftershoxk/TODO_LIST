package api

import (
	"net/http"
	"os"
	"todo_list/pkg/db"
)

func Init() {
	path := os.Getenv("TODO_DBFILE")
	if path == "" {
		path = "scheduler.db"
	}

	err := db.Init(path)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)

	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))

}
