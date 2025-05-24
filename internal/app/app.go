package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LikhithMar14/workout-tracker/internal/api"
	"github.com/LikhithMar14/workout-tracker/internal/store"
	"github.com/LikhithMar14/workout-tracker/migrations"
	"github.com/LikhithMar14/workout-tracker/pkg"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	Config         pkg.Config
	DB             *sql.DB
}

func NewApplication(cfg pkg.Config) (*Application, error) {
	pgDB, err := store.Open(cfg)
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	workoutStore := store.NewPostgressWorkoutStore(pgDB)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	userStore := store.NewPostgresUserStore(pgDB)
	userHandler := api.NewUserHandler(userStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		Config:         cfg,
		DB:             pgDB,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
