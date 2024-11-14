package repository

import (
	"database/sql"

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
