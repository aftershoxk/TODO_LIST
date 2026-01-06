package db

import (
	"database/sql"
	"os"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255),
	comment TEXT,
	repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
