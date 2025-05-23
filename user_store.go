package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string)(bool, error){
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))

	if err != nil{
		switch{
		case errors.Is(err,bcrypt.ErrMismatchedHashAndPassword):
			return false,nil
		default:
			return false,err
		}

	}

	return true,nil
}

var AnonymousUser = &User{}

func (u *User)IsAnonymous() bool {
	return u == AnonymousUser
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query :=
		`
			INSERT INTO users (username,email,password_hash,bio)
			VALUES ($1,$2,$3,$4)
			RETURNING id, created_at,updated_at

		`

	err := pg.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}
func(pg *PostgresUserStore)GetUserByUsername(username string)(*User,error){
	user := &User{
		PasswordHash: password{},
	}
	query := `
		SELECT id, username, bio, created_at,updated_at
		FROM users
		WHERE username=$1

	`

	err := pg.db.QueryRow(query,username).Scan(&user.ID,&user.Username,&user.PasswordHash.hash,&user.Bio,&user.CreatedAt,&user.UpdatedAt)

	if err != nil {
		switch{
		case errors.Is(err,sql.ErrNoRows):
			return nil,nil
		default:
			return nil,err
		}
	}

	return user, nil

}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, bio = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	result, err := s.db.Exec(query, user.Username, user.Email, user.Bio, user.ID)
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
