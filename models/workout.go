package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Workout struct {
	ID        int    `json:"id"`
	Steps     int    `json:"steps"`
	Calories  int    `json:"calories"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// Define the MySQL database connection string
const (
	DBHost  = "localhost"
	DBPort  = ":3306"
	DBUser  = "root"
	DBPass  = "password"
	DBDbase = "workouts"
)

// Handler for GET /workouts
func GetWorkoutsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, steps, calories, start_time, end_time FROM workouts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	workouts := make([]Workout, 0)
	for rows.Next() {
		var workout Workout
		if err := rows.Scan(&workout.ID, &workout.Steps, &workout.Calories, &workout.StartTime, &workout.EndTime); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

// Handler for GET /workouts/{id}
func GetWorkoutHandlerByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := r.URL.Query()
	id, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id, steps, calories, start_time, end_time FROM workouts WHERE id=?", id)
	var workout Workout
	if err := row.Scan(&workout.ID, &workout.Steps, &workout.Calories, &workout.StartTime, &workout.EndTime); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Workout with ID %d not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// Handler for POST /workouts
func CreateWorkoutHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Parse the request body to get the workout data
	var workout Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert the workout into the database
	result, err := db.ExecContext(r.Context(), `
		INSERT INTO workouts (steps, calories, start_time, end_time)
		VALUES (?, ?, ?, ?)
	`, workout.Steps, workout.Calories, workout.StartTime, workout.EndTime)
	if err != nil {
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly created workout
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get ID of created workout", http.StatusInternalServerError)
		return
	}

	// Set the ID field of the workout object and write it to the response
	workout.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(workout)
	if err != nil {
		http.Error(w, "Failed to encode workout JSON", http.StatusInternalServerError)
		return
	}
}

// Create a function to handle deleting workout records
func DeleteWorkout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parse the workout ID from the URL path
	id := r.URL.Path[len("/workouts/"):]

	// Delete the workout record from the database
	stmt, err := db.Prepare("DELETE FROM workouts WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	// Send a response back to the client
	fmt.Fprintf(w, "Deleted %d workout record(s)", rowsAffected)
}

// Create a function to handle updating workout records
func UpdateWorkout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parse the workout ID from the URL path
	id := r.URL.Path[len("/workouts/"):]

	// Parse the workout data from the request body
	var workout Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the workout record in the database
	stmt, err := db.Prepare("UPDATE workouts SET steps=?, calories=?, workout_start_time=?, workout_end_time=? WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(workout.Steps, workout.Calories, workout.StartTime, workout.EndTime, id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	// Send a response back to the client
	fmt.Fprintf(w, "Updated %d workout record(s)", rowsAffected)
}
