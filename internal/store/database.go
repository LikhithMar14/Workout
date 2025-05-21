package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/LikhithMar14/workout-tracker/pkg"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)



func Open(cfg pkg.Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
    )
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, fmt.Errorf("db: open %w", err)
    }
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("db: ping %w", err)
    }
    return db, nil
}
func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	//we use anon Fn when we have so many things to close , in this case we can also use normally
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil{
		return fmt.Errorf("migrate: %w",err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}