package repository

import (
	"database/sql"
	"errors"

	"github.com/mauFade/playzy/internal/model"
)

type UserRepositoryInterface interface {
	FindByID(id string) (*model.UserModel, error)
	FindByEmail(email string) (*model.UserModel, error)
	FindByPhone(phone string) (*model.UserModel, error)
	FindByGamertag(gamertag string) (*model.UserModel, error)
	Create(user *model.UserModel) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(d *sql.DB) *UserRepository {
	r := &UserRepository{
		db: d,
	}

	r.db.Exec("CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY, name VARCHAR NOT NULL, email VARCHAR NOT NULL, phone VARCHAR NOT NULL, password VARCHAR NOT NULL, gamertag VARCHAR NOT NULL, is_deleted BOOLEAN NOT NULL, deleted_at TIMESTAMP NULL, updated_at TIMESTAMP NOT NULL, created_at TIMESTAMP NOT NULL)")

	return r
}

func (r *UserRepository) Create(user *model.UserModel) error {
	query := `INSERT INTO users 
	(id, name, email, phone, gamertag, password, is_deleted, deleted_at, updated_at, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, 'false', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

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

func (r *UserRepository) FindByID(id string) (*model.UserModel, error) {
	query := "SELECT * FROM users WHERE id = $1"
	row := r.db.QueryRow(query, id)

	var user model.UserModel

	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Gamertag,
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

func (r *UserRepository) FindByEmail(email string) (*model.UserModel, error) {
	query := "SELECT * FROM users WHERE email = $1"
	row := r.db.QueryRow(query, email)

	var user model.UserModel

	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Gamertag,
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

func (r *UserRepository) FindByGamertag(gamertag string) (*model.UserModel, error) {
	query := "SELECT * FROM users WHERE gamertag = $1"
	row := r.db.QueryRow(query, gamertag)

	var user model.UserModel

	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Gamertag,
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

func (r *UserRepository) FindByPhone(phone string) (*model.UserModel, error) {
	query := "SELECT * FROM users WHERE phone = $1"
	row := r.db.QueryRow(query, phone)

	var user model.UserModel

	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Gamertag,
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
