package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LikhithMar14/workout-tracker/internal/api"
)

type Application struct{
	Logger *log.Logger
	WorkoutHandler *api.WorkoutHandler
}


func NewApplication() (*Application, error){
	logger := log.New(os.Stdout,"",log.Ldate|log.Ltime|log.Lshortfile)
	workoutHandler := api.NewWorkoutHandler(logger)
	app := &Application{
		Logger: logger,
		WorkoutHandler: workoutHandler,
	}
	return app,nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Status is available\n")
}