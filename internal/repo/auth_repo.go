package repo

import (
	"fmt"
	"server/internal/model"

	"github.com/jackc/pgx"
)


type AuthRepo struct {
	db *pgx.Conn
}

func NewAuthRepo(db *pgx.Conn) *AuthRepo {
	return  &AuthRepo{db: db}
}

// CreateUser queries db to create a new user
func (r *AuthRepo) CreateUser(u *model.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	row := r.db.QueryRow(query, u.Username, u.Email, u.PasswordHash)
	scanErr := row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if  scanErr != nil {
	  return fmt.Errorf("create user scan returning: %w", scanErr)
	} else {
		return nil
	}
}

// GetByEmail uses db connection to query users table by username
func (r *AuthRepo) GetByEmail(email string) (*model.User, error){
	u := new(model.User)
	query := `SELECT id, username, email, password_hash FROM users WHERE email=$1`
	row := r.db.QueryRow(query, email)
	
	scanErr := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash)
	
	if scanErr != nil {
		return nil, scanErr
	} else {
		return u, nil
	}
}