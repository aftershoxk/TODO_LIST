package main

import (
	"log"
	"net/http"
	"os"
	"todo_list/pkg/api"
	"todo_list/pkg/db"

	_ "modernc.org/sqlite"
)

// Функция для получения порта из переменной окружения
func GetPort() string {
	port := ":"
	data := os.Getenv("TODO_PORT")
	if data == "" {
		port += "7540"
	} else {
		port += data
	}
	return port
}

// Функция для получения пути к базе данных из переменной окружения
func DBPath() string {
	data := os.Getenv("TODO_DBFILE")
	if data == "" {
		return "scheduler.db"
	}
	return data
}

func main() {
	// инициализация БД
	path := DBPath()
	err := db.Init(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("error closing database:", err)
		}
	}()

	// файл-сервер
	fileHandler := http.FileServer(http.Dir("./web"))
	http.Handle("/", fileHandler)

	//Регистрация HTTP-обработчика
	api.Init()

	// запуск сервера
	port := GetPort()
	log.Printf("server starting on %s", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
