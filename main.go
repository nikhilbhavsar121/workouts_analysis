package main

import (
	"aubergine_test/routes"
	"aubergine_test/tasks"
	"aubergine_test/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Workout struct {
	ID        int    `json:"id"`
	Steps     int    `json:"steps"`
	Calories  int    `json:"calories"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func main() {
	LoadConfig()

	conn, dbErr := utils.GetDBConnection()

	if dbErr != nil {
		fmt.Println("Cannot connect to database", dbErr)
		return
	}

	defer conn.Close()

	routes.SetupRoutes(conn)
	port := os.Getenv("PORT")
	startWorkers()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}
}

func startWorkers() {
	go tasks.DailyAggregationsRunEvery(time.Duration(1) * time.Hour)

	go tasks.WeeklyAggregationsRunEvery(time.Duration(168) * time.Hour)

}
