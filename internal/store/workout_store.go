package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Workout struct {
	ID              int            `json:"id"`
	UserID          int            `json:"user_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description,omitempty"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  *int           `json:"calories_burned,omitempty"`
	Entries         []WorkoutEntry `json:"entries,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type WorkoutEntry struct {
	ID              int       `json:"id"`
	ExerciseName    string    `json:"exercise_name"`
	Sets            int       `json:"sets"`
	Reps            *int      `json:"reps,omitempty"`
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
	Weight          *float64  `json:"weight,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
	OrderIndex      int       `json:"order_index"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PostgressWorkoutStore struct {
	db *sql.DB
}

func NewPostgressWorkoutStore(db *sql.DB) *PostgressWorkoutStore {
	return &PostgressWorkoutStore{db}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(int64) (*Workout, error)
	GetWorkoutByIDAndUserID(workoutID int64, userID int) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkoutByID(int64) error
	DeleteWorkoutByIDAndUserID(workoutID int64, userID int) error
}

func (pg *PostgressWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(query, workout.UserID, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID, &workout.CreatedAt, &workout.UpdatedAt)
	if err != nil {
		return nil, err
	}

	for i := range workout.Entries {
		entry := &workout.Entries[i]
		query := `
			INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgressWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	workout := &Workout{}
	query := `
	SELECT id, user_id, title, description, duration_minutes, calories_burned, created_at, updated_at
	FROM workouts
	WHERE id=$1
	`
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.UserID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned, &workout.CreatedAt, &workout.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	entryQuery := `
  SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
  FROM workout_entries
  WHERE workout_id = $1
  ORDER BY order_index
  `

	rows, err := pg.db.Query(entryQuery, id)
	fmt.Print(err)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry
		err = rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (pg *PostgressWorkoutStore) GetWorkoutByIDAndUserID(workoutID int64, userID int) (*Workout, error) {
	workout := &Workout{}
	query := `
	SELECT id, user_id, title, description, duration_minutes, calories_burned, created_at, updated_at
	FROM workouts
	WHERE id=$1 AND user_id=$2
	`
	err := pg.db.QueryRow(query, workoutID, userID).Scan(&workout.ID, &workout.UserID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned, &workout.CreatedAt, &workout.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	entryQuery := `
  SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index, created_at, updated_at
  FROM workout_entries
  WHERE workout_id = $1
  ORDER BY order_index
  `

	rows, err := pg.db.Query(entryQuery, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry
		err = rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (pg *PostgressWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query :=
		`
		UPDATE workouts SET title=$1, description=$2, duration_minutes=$3, calories_burned=$4
		WHERE id=$5
	
	`
	//we use exec when we are doing put/patch/delete or when we are not returning anything
	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	_, err = tx.Exec(`DELETE FROM workout_entries WHERE workout_id = $1`, workout.ID)

	if err != nil {
		return err
	}
	for _, entry := range workout.Entries {
		query := `
		INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		_, err := tx.Exec(query,
			workout.ID,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.DurationSeconds,
			entry.Weight,
			entry.Notes,
			entry.OrderIndex,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()

}

func (pg *PostgressWorkoutStore) DeleteWorkoutByID(id int64) error {
	query := `
  DELETE from workouts
  WHERE id = $1
  `

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgressWorkoutStore) DeleteWorkoutByIDAndUserID(workoutID int64, userID int) error {
	query := `
  DELETE from workouts
  WHERE id = $1 AND user_id = $2
  `

	result, err := pg.db.Exec(query, workoutID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	if err != nil {
		return err
	}

	return nil
}
