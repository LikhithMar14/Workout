package routes

import (
	"net/http"

	"github.com/LikhithMar14/workout-tracker/internal/app"
	"github.com/go-chi/chi/v5"
)
func SetupRoutes(app *app.Application) http.Handler {
	r := chi.NewRouter()
	r.Get("/workouts/{id}",app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}",app.WorkoutHandler.HandleUpdateWorkoutByID)
	r.Delete("/workouts/{id}",app.WorkoutHandler.HandleDeleteWorkoutByID)
	r.Get("/health", app.HealthCheck)
	r.Post("/users",app.UserHandler.HandleRegisterUser)

	return r
}