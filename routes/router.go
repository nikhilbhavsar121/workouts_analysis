package routes

import (
	"aubergine_test/models"
	"database/sql"
	"net/http"
)

func SetupRoutes(db *sql.DB) {
	//TODO: send id as path param

	http.HandleFunc("/workouts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			models.GetWorkoutsHandler(w, r, db)
		case http.MethodPost:
			models.CreateWorkoutHandler(w, r, db)
		case http.MethodPut:
			models.UpdateWorkout(w, r, db)
		case http.MethodDelete:
			models.DeleteWorkout(w, r, db)
		default:
			http.NotFound(w, r)
		}
	})
	//TODO: send id as path param
	http.HandleFunc("/getworkouts", func(w http.ResponseWriter, r *http.Request) {
		models.GetWorkoutHandlerByID(w, r, db)
	})
}
