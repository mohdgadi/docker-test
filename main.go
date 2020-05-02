package main

import (
	"database/sql"
	"docker-test/todo"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize a connextion to DB. Sourcename contains the auth information, host information
	// obtained from the env.
	db, err := sql.Open("mysql", sourceName)
	if err != nil {
		panic(err)
	}

	Init(db)
}

// Init initializes the dependecies and boots up the server.
func Init(db *sql.DB) {
	// Initialize our model
	repo := todo.NewListRepository(db)
	err := repo.Init()
	if err != nil {
		panic(err)
	}

	// Pass our model to our service which will handle buisiness logic
	listService := todo.NewListService(repo)

	// start server
	http.HandleFunc("/save", listService.AddItem)
	http.HandleFunc("/delete", listService.DeleteItem)
	http.HandleFunc("/find", listService.FindItem)
	http.ListenAndServe(":8080", nil)
}
