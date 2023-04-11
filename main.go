package main

import (
	"aubergine_test/routes"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Workout struct {
	ID        int    `json:"id"`
	Steps     int    `json:"steps"`
	Calories  int    `json:"calories"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:Nikhil58@@tcp(127.0.0.1:3306)/workoutsdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	routes.SetupRoutes(db)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
