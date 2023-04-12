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

	var workouts []Workout
	err := json.NewDecoder(r.Body).Decode(&workouts)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, workout := range workouts {

		result, err := db.ExecContext(r.Context(), `
			INSERT INTO workouts (steps, calories, start_time, end_time)
			VALUES (?, ?, ?, ?)
		`, workout.Steps, workout.Calories, workout.StartTime, workout.EndTime)
		if err != nil {
			http.Error(w, "Failed to create workout", http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to get ID of created workout", http.StatusInternalServerError)
			return
		}

		workout.ID = int(id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(workout)
		if err != nil {
			http.Error(w, "Failed to encode workout JSON", http.StatusInternalServerError)
			return
		}
	}
}

// Create a function to handle deleting workout records
func DeleteWorkout(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := r.URL.Query()
	id, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
	vars := r.URL.Query()
	id, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Parse the workout data from the request body
	var workout Workout
	err = json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the workout record in the database
	stmt, err := db.Prepare("UPDATE workouts SET steps=?, calories=?, start_time=?, end_time=? WHERE id=?")
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
	// dateStr := workout.StartTime
	// layout := "2006-01-02 15:04:05"
	// dt, _ := time.Parse(layout, dateStr)
	// startDtStr := dt.Format("2006-01-02 15:00:00")
	// endDtStr := dt.Format("2006-01-02 15:59:00")
	// fmt.Println(startDtStr) // Output: 2023-04-07 22:00:00
	// fmt.Println(endDtStr)   // Output: 2023-04-07 22:59:00

	// Send a response back to the client
	fmt.Fprintf(w, "Updated %d workout record(s)", rowsAffected)
}
