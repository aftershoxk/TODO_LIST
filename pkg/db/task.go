package db

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

const DateFormat = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const selectTaskFields = `
	id,
	date,
	title,
	comment,
	repeat
`

func Tasks(limit int, search string) ([]*Task, error) {
	result := []*Task{}
	var rows *sql.Rows
	var err error

	if search == "" {
		query := `
			SELECT ` + selectTaskFields + `
			FROM scheduler
			ORDER BY date
			LIMIT ?
		`
		rows, err = db.Query(query, limit)
		if err != nil {
			return nil, err
		}
	} else {
		p, err := time.Parse("02.01.2006", search)
		if err == nil {
			date := p.Format(DateFormat)
			query := `
				SELECT ` + selectTaskFields + `
				FROM scheduler
				WHERE date = ?
				ORDER BY date
				LIMIT ?
			`
			rows, err = db.Query(query, date, limit)
			if err != nil {
				return nil, err
			}
		} else {
			pattern := "%" + search + "%"
			query := `
				SELECT ` + selectTaskFields + `
				FROM scheduler
				WHERE title LIKE ? OR comment LIKE ?
				ORDER BY date
				LIMIT ?
			`
			rows, err = db.Query(query, pattern, pattern, limit)
			if err != nil {
				return nil, err
			}
		}
	}

	defer rows.Close()

	for rows.Next() {
		var task Task
		var idInt int64

		err = rows.Scan(
			&idInt,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)
		if err != nil {
			return nil, err
		}

		task.ID = strconv.Itoa(int(idInt))
		result = append(result, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`

	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	var task Task
	var idInt int64

	query := `
		SELECT ` + selectTaskFields + `
		FROM scheduler
		WHERE id = ?
	`

	err := db.QueryRow(query, id).Scan(
		&idInt,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &task, errors.New("task doesn`t exist")
		}
		return &task, err
	}

	task.ID = strconv.Itoa(int(idInt))
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	res, err := db.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("incorrect id")
	}

	return nil
}

func DeleteTask(id string) error {
	if id == "" {
		return errors.New("empty id")
	}

	res, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("incorrect id")
	}

	return nil
}

func UpdateDate(id string, next string) error {
	if id == "" || next == "" {
		return errors.New("invalid id or date")
	}

	res, err := db.Exec(
		"UPDATE scheduler SET date = ? WHERE id = ?",
		next,
		id,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("incorrect id")
	}

	return nil
}
