package tasks

import (
	"aubergine_test/models"
	"aubergine_test/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func WeeklyAggregationsRunEvery(d time.Duration) {
	conn, dbErr := utils.GetDBConnection()

	if dbErr != nil {
		fmt.Println("Cannot connect to database", dbErr)
		return
	}

	defer conn.Close()
	for range time.Tick(d) {
		log.Println("task# weekly_aggregations task is started")

		now := time.Now()
		lastWeekStart := now.AddDate(0, 0, -7).Truncate(24 * time.Hour)
		lastWeekEnd := lastWeekStart.AddDate(0, 0, 7)
		log.Println("task# find the workout of last hour", lastWeekStart, lastWeekEnd)

		calculateAndSubmitBasedOnRange(lastWeekStart, lastWeekEnd, conn)
		log.Println("task# weekly_aggregations task is finished at", time.Now(), "Next will start at", time.Now().AddDate(0, 0, 7))
	}
}

func DailyAggregationsRunEvery(d time.Duration) {
	conn, dbErr := utils.GetDBConnection()

	if dbErr != nil {
		fmt.Println("Cannot connect to database", dbErr)
		return
	}

	defer conn.Close()
	for range time.Tick(d) {
		log.Println("task# daily_aggregations task is started")
		now := time.Now()
		// If it's a new hour
		// Get the start and end times for this hour
		hourEnd := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		hourStart := hourEnd.Add(-1 * time.Hour)
		log.Println("Find the workout of last hour", hourStart, hourEnd)

		calculateAndSubmitBasedOnRange(hourStart, hourEnd, conn)
		log.Println("task# daily_aggregations task is finished at", time.Now(), "Next will start at", time.Now().Add(1*time.Hour))
	}
}

func calculateAndSubmitBasedOnRange(hourStart, hourEnd time.Time, conn *sql.DB) {
	// Query the workouts table for all workouts within this hour
	workouts, err := getWorkoutsForRange(hourStart, hourEnd, conn)
	if err != nil {
		return
	}
	if len(workouts) == 0 {
		log.Println("task# no workout found for starttime : ", hourStart, "endtime :", hourEnd)
		return
	}
	// Calculate the total steps and calories for the hour
	totalSteps := 0
	totalCalories := 0
	for _, workout := range workouts {
		totalSteps += workout.Steps
		totalCalories += workout.Calories
	}

	// Insert a new row into the daily_aggregations table for this hour
	insertEntry(hourStart, hourEnd, totalSteps, totalCalories, conn)
}

func insertEntry(hourStart, hourEnd time.Time, totalSteps int, totalCalories int, conn *sql.DB) {
	result, err := conn.ExecContext(context.Background(), `
	INSERT INTO workouts (steps, calories, start_time, end_time)
	VALUES (?, ?, ?, ?)
	`, totalSteps, totalCalories, hourStart, hourEnd)

	if err != nil {
		log.Printf("Failed to insert workout: %v\n", err)
		return
	}
	// Get the ID of the newly created workout
	id, err := result.LastInsertId()
	if err != nil {
		return
	}
	log.Println("task# result created with id", id)
}

func getWorkoutsForRange(hourStart time.Time, hourEnd time.Time, conn *sql.DB) ([]models.Workout, error) {

	rows, err := conn.Query("SELECT id, steps, calories, start_time, end_time FROM workouts")
	if err != nil {
		fmt.Println("Cannot connect to database", err)
		return nil, err
	}
	defer rows.Close()

	workouts := make([]models.Workout, 0)
	for rows.Next() {
		var workout models.Workout
		if err := rows.Scan(&workout.ID, &workout.Steps, &workout.Calories, &workout.StartTime, &workout.EndTime); err != nil {
			fmt.Println("Cannot connect to database", err)
			return nil, err
		}
		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Cannot connect to database", err)
		return nil, err
	}
	return nil, err

}
