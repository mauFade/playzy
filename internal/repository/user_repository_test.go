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
	// Criando o mock do banco
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	// Inicializando o repositório com o banco mockado
	repo := repository.NewUserRepository(db)

	// Preparando os dados de entrada
	id := uuid.New()
	user := &model.UserModel{
		ID:        id,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
		Gamertag:  "johnny",
		Password:  "password123",
		Deleted:   false,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	// Esperando que a query INSERT seja chamada
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.GetID(), user.GetName(), user.GetEmail(), user.GetPhone(), user.GetGamertag(), user.GetPassword()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Chamando o método Create
	err = repo.Create(user)

	// Verificando se a query foi executada corretamente e se o erro é nil
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByEmail(t *testing.T) {
	// Criando o mock do banco
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	// Inicializando o repositório com o banco mockado
	repo := repository.NewUserRepository(db)

	// Preparando os dados de entrada
	email := "john.doe@example.com"
	user := &model.UserModel{
		ID:        uuid.New(),
		Name:      "John Doe",
		Email:     email,
		Phone:     "1234567890",
		Gamertag:  "johnny",
		Password:  "password123",
		Deleted:   false,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	// Esperando que a query SELECT seja chamada
	mock.ExpectQuery("SELECT \\* FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone", "gamertag", "password", "is_deleted", "deleted_at", "updated_at", "created_at"}).
			AddRow(user.GetID(), user.GetName(), user.GetEmail(), user.GetPhone(), user.GetGamertag(), user.GetPassword(), user.IsDeleted(), nil, user.GetUpdatedAt(), user.GetCreatedAt()))

	// Chamando o método FindByEmail
	result, err := repo.FindByEmail(email)

	// Verificando se a query foi executada corretamente e se o resultado é o esperado
	assert.NoError(t, err)
	assert.Equal(t, user.GetEmail(), result.GetEmail())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByPhone(t *testing.T) {
	// Criando o mock do banco
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	// Inicializando o repositório com o banco mockado
	repo := repository.NewUserRepository(db)

	// Preparando os dados de entrada
	phone := "1234567890"
	user := &model.UserModel{
		ID:        uuid.New(),
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Phone:     phone,
		Gamertag:  "johnny",
		Password:  "password123",
		Deleted:   false,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	// Esperando que a query SELECT seja chamada
	mock.ExpectQuery("SELECT \\* FROM users WHERE phone = ?").
		WithArgs(phone).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone", "gamertag", "password", "is_deleted", "deleted_at", "updated_at", "created_at"}).
			AddRow(user.GetID(), user.GetName(), user.GetEmail(), user.GetPhone(), user.GetGamertag(), user.GetPassword(), user.IsDeleted(), nil, user.GetUpdatedAt(), user.GetCreatedAt()))

	// Chamando o método FindByPhone
	result, err := repo.FindByPhone(phone)

	// Verificando se a query foi executada corretamente e se o resultado é o esperado
	assert.NoError(t, err)
	assert.Equal(t, user.GetPhone(), result.GetPhone())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByGamertag(t *testing.T) {
	// Criando o mock do banco
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	// Inicializando o repositório com o banco mockado
	repo := repository.NewUserRepository(db)

	// Preparando os dados de entrada
	gamertag := "johnny"
	user := &model.UserModel{
		ID:        uuid.New(),
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
		Gamertag:  gamertag,
		Password:  "password123",
		Deleted:   false,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	// Esperando que a query SELECT seja chamada
	mock.ExpectQuery("SELECT \\* FROM users WHERE gamertag = ?").
		WithArgs(gamertag).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone", "gamertag", "password", "is_deleted", "deleted_at", "updated_at", "created_at"}).
			AddRow(user.GetID(), user.GetName(), user.GetEmail(), user.GetPhone(), user.GetGamertag(), user.GetPassword(), user.IsDeleted(), nil, user.GetUpdatedAt(), user.GetCreatedAt()))

	// Chamando o método FindByGamertag
	result, err := repo.FindByGamertag(gamertag)

	// Verificando se a query foi executada corretamente e se o resultado é o esperado
	assert.NoError(t, err)
	assert.Equal(t, user.GetGamertag(), result.GetGamertag())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByEmailNotFound(t *testing.T) {
	// Criando o mock do banco
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar o mock do banco: %v", err)
	}
	defer db.Close()

	// Inicializando o repositório com o banco mockado
	repo := repository.NewUserRepository(db)

	// Preparando o dado de entrada
	email := "nonexistent@example.com"

	// Esperando que a query SELECT seja chamada e que não retorne nenhum resultado
	mock.ExpectQuery("SELECT \\* FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	// Chamando o método FindByEmail
	result, err := repo.FindByEmail(email)

	// Verificando se o erro retornado é nil e o resultado também é nil
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
