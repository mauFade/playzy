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
	query := `INSERT INTO users values 
	(id, name, email, phone, gamertag, password, is_deleted, deleted_at, updated_at, created_at) 
	VALUES (?, ?, ?, ?, ?, ?, 'false', NULL, NOW(), NOW())`

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
		user.GetID(),
		user.GetName(),
		user.GetEmail(),
		user.GetPhone(),
		user.GetGamertag(),
		user.GetPassword(),
		user.IsDeleted(),
		user.GetDeletedAt(),
		user.GetUpdatedAt(),
		user.GetCreatedAt(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
