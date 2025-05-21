package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/LikhithMar14/workout-tracker/internal/app"
	"github.com/LikhithMar14/workout-tracker/internal/routes"
	"github.com/LikhithMar14/workout-tracker/pkg"
)

func main() {

    var (
        port     int
        host     string
        user     string
        password string
        dbname   string
        dbPort   int
        sslmode  string
    )

    flag.IntVar(&port, "port", 8080, "Go backend server port")
    flag.StringVar(&host, "host", "localhost", "Host of Database")
    flag.StringVar(&user, "user", "postgres", "User of Database")
    flag.StringVar(&password, "password", "postgres", "Password of Database")
    flag.StringVar(&dbname, "dbname", "postgres", "Database name")
    flag.IntVar(&dbPort, "dbport", 5432, "Database port")
    flag.StringVar(&sslmode, "sslmode", "disable", "SSL mode for database connection")

    flag.Parse()

    cfg := pkg.Config{
        DBHost:     host,
        DBUser:     user,
        DBPassword: password,
        DBName:     dbname,
        DBPort:     dbPort,
        DBSSLMode:  sslmode,
        ServerPort: port,
    }

    app, err := app.NewApplication(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to initialize application: %w", err))
	}

	r := routes.SetupRoutes(app)


	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("Starting server on port %d...\n", port)


	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		app.Logger.Fatalf("Server failed: %v", err)
	}
}
