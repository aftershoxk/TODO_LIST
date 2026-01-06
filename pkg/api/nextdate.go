package api

import (
	"net/http"
	"time"
	n "todo_list/pkg/repeat"
)

const DateFormat = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")
	var nowTime time.Time

	if dateStr == "" || repeatStr == "" {
		writeError(w, http.StatusBadRequest, "missing date or repeat")
		return
	}
	if nowStr == "" {
		nowTime = time.Now()
	} else {
		var err error
		nowTime, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid now format")
			return
		}
	}

	result, err := n.NextDate(nowTime, dateStr, repeatStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
