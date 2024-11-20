package repository

import (
	"database/sql"
	"errors"

	"github.com/mauFade/playzy/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(d *sql.DB) *UserRepository {
	return &UserRepository{
		db: d,
	}
}

func (r *UserRepository) Create(user *model.UserModel) error {
	query := `INSERT INTO users 
	(id, name, email, phone, gamertag, password, is_deleted, deleted_at, updated_at, created_at) 
	VALUES (?, ?, ?, ?, ?, ?, 'false', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := r.db.Exec(query,
		user.GetID(),
		user.GetName(),
		user.GetEmail(),
		user.GetPhone(),
		user.GetGamertag(),
		user.GetPassword(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.UserModel, error) {
	query := "SELECT * FROM users WHERE email = ?"
	row := r.db.QueryRow(query, email)

	var user model.UserModel

	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Gamertag,
		&user.Password,
		&user.Deleted,
		&user.DeletedAt,
		&user.UpdatedAt,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
