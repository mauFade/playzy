package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/mauFade/playzy/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	id := uuid.New()
	name := "John Doe"
	email := "john.doe@example.com"
	phone := "1234567890"
	gamertag := "johnny"
	password := "password123"
	deleted := false
	var deletedAt *time.Time = nil
	updatedAt := time.Now()
	createdAt := time.Now()

	user := model.NewUserModel(id, name, email, phone, gamertag, password, deleted, deletedAt, updatedAt, createdAt)

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.GetID(), user.GetName(), user.GetEmail(), user.GetPhone(), user.GetGamertag(), user.GetPassword()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	id := uuid.New()
	name := "John Doe"
	email := "john.doe@example.com"
	phone := "1234567890"
	gamertag := "johnny"
	password := "password123"
	deleted := false
	var deletedAt *time.Time = nil
	updatedAt := time.Now()
	createdAt := time.Now()

	user := model.NewUserModel(id, name, email, phone, gamertag, password, deleted, deletedAt, updatedAt, createdAt)

	mock.ExpectQuery("SELECT \\* FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone", "gamertag", "password", "is_deleted", "deleted_at", "updated_at", "created_at"}).
			AddRow(user.GetID(), user.GetName(), user.GetEmail(), user.GetPhone(), user.GetGamertag(), user.GetPassword(), user.IsDeleted(), nil, user.GetUpdatedAt(), user.GetCreatedAt()))

	result, err := repo.FindByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, user.GetEmail(), result.GetEmail())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByPhone(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	id := uuid.New()
	name := "John Doe"
	email := "john.doe@example.com"
	phone := "1234567890"
	gamertag := "johnny"
	password := "password123"
	deleted := false
	var deletedAt *time.Time = nil
	updatedAt := time.Now()
	createdAt := time.Now()

	user := model.NewUserModel(id, name, email, phone, gamertag, password, deleted, deletedAt, updatedAt, createdAt)

	mock.ExpectQuery("SELECT \\* FROM users WHERE phone = ?").
		WithArgs(phone).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone", "gamertag", "password", "is_deleted", "deleted_at", "updated_at", "created_at"}).
			AddRow(user.GetID(), user.GetName(), user.GetEmail(), user.GetPhone(), user.GetGamertag(), user.GetPassword(), user.IsDeleted(), nil, user.GetUpdatedAt(), user.GetCreatedAt()))

	result, err := repo.FindByPhone(phone)

	assert.NoError(t, err)
	assert.Equal(t, user.GetPhone(), result.GetPhone())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByGamertag(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	id := uuid.New()
	name := "John Doe"
	email := "john.doe@example.com"
	phone := "1234567890"
	gamertag := "johnny"
	password := "password123"
	deleted := false
	var deletedAt *time.Time = nil
	updatedAt := time.Now()
	createdAt := time.Now()

	user := model.NewUserModel(id, name, email, phone, gamertag, password, deleted, deletedAt, updatedAt, createdAt)

	mock.ExpectQuery("SELECT \\* FROM users WHERE gamertag = ?").
		WithArgs(gamertag).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone", "gamertag", "password", "is_deleted", "deleted_at", "updated_at", "created_at"}).
			AddRow(user.GetID(), user.GetName(), user.GetEmail(), user.GetPhone(), user.GetGamertag(), user.GetPassword(), user.IsDeleted(), nil, user.GetUpdatedAt(), user.GetCreatedAt()))

	result, err := repo.FindByGamertag(gamertag)

	assert.NoError(t, err)
	assert.Equal(t, user.GetGamertag(), result.GetGamertag())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByEmailNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	email := "nonexistent@example.com"

	mock.ExpectQuery("SELECT \\* FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.FindByEmail(email)

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
