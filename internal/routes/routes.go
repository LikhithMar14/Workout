package routes

import (
	"net/http"
	"time"

	"github.com/LikhithMar14/workout-tracker/internal/app"
	"github.com/LikhithMar14/workout-tracker/internal/auth"
	"github.com/LikhithMar14/workout-tracker/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) http.Handler {
	r := chi.NewRouter()

	// Create JWT authenticator (in production, get these from environment variables)
	authenticator := auth.NewJWTAuthenticator("mysecret", "workout-tracker-users", "workout-tracker-app")

	// Create middleware instance
	mw := middleware.NewMiddleware(app.Logger, authenticator)

	// Create rate limiter (100 requests per minute)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	// Global middleware (applied to all routes)
	r.Use(mw.RecoverPanic)
	r.Use(mw.CORS)
	r.Use(mw.RequestLogger)
	r.Use(mw.ContentType)
	r.Use(mw.RateLimit(rateLimiter))

	// Public routes (no authentication required)
	r.Post("/register", app.UserHandler.HandleRegisterUser)
	r.Post("/login", app.UserHandler.HandleLoginUser)
	r.Get("/health", app.HealthCheck)

	// Protected routes (authentication required)
	r.Group(func(r chi.Router) {
		r.Use(mw.RequireAuth)

		// Workout routes
		r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
		r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
		r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
		r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)
	})

	return r
}
