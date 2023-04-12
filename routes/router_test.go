package routes

import (
	"aubergine_test/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSetupRoutes(t *testing.T) {

	conn, err := sql.Open("mysql", "root:Nikhil58@@tcp(127.0.0.1:3306)/workoutdb")
	if err != nil {
		fmt.Println("Cannot connect to database", err)
		return
	}

	// create a new workout
	workoutData := []byte(`[{
        "steps": 80,
        "calories": 0,
        "start_time": "2023-04-07 22:41:00",
        "end_time": "2023-04-07 22:41:59"
    }]`)

	req, err := http.NewRequest("POST", "/workouts", bytes.NewBuffer(workoutData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		models.CreateWorkoutHandler(w, r, conn)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// get all workouts
	req, err = http.NewRequest("GET", "/workouts", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		models.GetWorkoutsHandler(w, r, conn)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// parse the response body to get the ID of the created workout
	var workouts []models.Workout
	err = json.Unmarshal(rr.Body.Bytes(), &workouts)
	if err != nil {
		t.Fatal(err)
	}

	workoutID := workouts[len(workouts)-1].ID

	// update the created workout
	updatedWorkoutData := []byte(`{
        "steps": 100,
        "calories": 0,
        "start_time": "2023-04-07 22:41:00",
        "end_time": "2023-04-07 22:41:59"
    }`)
	req, err = http.NewRequest("PUT", fmt.Sprintf("/workouts/?id=%d", workoutID), bytes.NewBuffer(updatedWorkoutData))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		models.UpdateWorkout(w, r, conn)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// delete the created workout
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/workouts/?id=%d", workoutID), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		models.DeleteWorkout(w, r, conn)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
