package routes

import (
	"aubergine_test/models"
	"database/sql"
	"net/http"
)

func SetupRoutes(db *sql.DB) {
	http.HandleFunc("/workouts", func(w http.ResponseWriter, r *http.Request) {
		models.GetWorkoutsHandler(w, r, db)
	})
	http.HandleFunc("/workouts/", func(w http.ResponseWriter, r *http.Request) {
		models.GetWorkoutHandlerByID(w, r, db)
	})
	http.HandleFunc("/workouts", func(w http.ResponseWriter, r *http.Request) {
		models.CreateWorkoutHandler(w, r, db)
	})
	http.HandleFunc("/workouts/", func(w http.ResponseWriter, r *http.Request) {
		models.UpdateWorkout(w, r, db)
	})
	http.HandleFunc("/workouts/", func(w http.ResponseWriter, r *http.Request) {
		models.DeleteWorkout(w, r, db)
	})

}
