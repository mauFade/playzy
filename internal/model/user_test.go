package model_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewUserModel(t *testing.T) {
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

	assert.Equal(t, id, user.GetID())
	assert.Equal(t, name, user.GetName())
	assert.Equal(t, email, user.GetEmail())
	assert.Equal(t, phone, user.GetPhone())
	assert.Equal(t, gamertag, user.GetGamertag())
	assert.Equal(t, password, user.GetPassword())
	assert.Equal(t, deleted, user.IsDeleted())
	assert.Nil(t, user.GetDeletedAt())
	assert.WithinDuration(t, updatedAt, user.GetUpdatedAt(), time.Second)
	assert.WithinDuration(t, createdAt, user.GetCreatedAt(), time.Second)
}

func TestUserSettersAndGetters(t *testing.T) {
	user := &model.UserModel{}

	user.SetName("Jane Doe")
	user.SetEmail("jane.doe@example.com")
	user.SetPhone("9876543210")
	user.SetGamertag("janedoe")
	user.SetPassword("password456")
	user.SetDeleted(true)

	assert.Equal(t, "Jane Doe", user.GetName())
	assert.Equal(t, "jane.doe@example.com", user.GetEmail())
	assert.Equal(t, "9876543210", user.GetPhone())
	assert.Equal(t, "janedoe", user.GetGamertag())
	assert.Equal(t, "password456", user.GetPassword())
	assert.True(t, user.IsDeleted())
}
