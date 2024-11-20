package model_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewUserModel(t *testing.T) {
	// Prepare test data
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

	// Create a new UserModel using the constructor
	user := model.NewUserModel(id, name, email, phone, gamertag, password, deleted, deletedAt, updatedAt, createdAt)

	// Validate that the returned user model has the correct values
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

func TestSettersAndGetters(t *testing.T) {
	// Prepare initial user model
	user := &model.UserModel{}

	// Test Setters
	user.SetName("Jane Doe")
	user.SetEmail("jane.doe@example.com")
	user.SetPhone("9876543210")
	user.SetGamertag("janedoe")
	user.SetPassword("password456")
	user.SetDeleted(true)

	// Test Getters
	assert.Equal(t, "Jane Doe", user.GetName())
	assert.Equal(t, "jane.doe@example.com", user.GetEmail())
	assert.Equal(t, "9876543210", user.GetPhone())
	assert.Equal(t, "janedoe", user.GetGamertag())
	assert.Equal(t, "password456", user.GetPassword())
	assert.True(t, user.IsDeleted())
}

func TestDeletedAt(t *testing.T) {
	// Prepare a non-deleted user with no deletedAt
	user := &model.UserModel{
		Deleted:   false,
		DeletedAt: nil,
	}

	// Test GetDeletedAt for a non-deleted user
	assert.Nil(t, user.GetDeletedAt())

	// Set the user as deleted and assign a deletedAt time
	deletedAt := time.Now()
	user.SetDeleted(true)
	user.DeletedAt = &deletedAt

	// Test GetDeletedAt for a deleted user
	assert.NotNil(t, user.GetDeletedAt())
	assert.WithinDuration(t, deletedAt, *user.GetDeletedAt(), time.Second)
}

func TestUpdateTimes(t *testing.T) {
	// Prepare user model with initial timestamps
	initialUpdatedAt := time.Now().Add(-time.Hour)
	initialCreatedAt := time.Now().Add(-2 * time.Hour)

	user := &model.UserModel{
		UpdatedAt: initialUpdatedAt,
		CreatedAt: initialCreatedAt,
	}

	// Test initial values
	assert.WithinDuration(t, initialUpdatedAt, user.GetUpdatedAt(), time.Second)
	assert.WithinDuration(t, initialCreatedAt, user.GetCreatedAt(), time.Second)

	// Update user timestamps
	newUpdatedAt := time.Now()
	newCreatedAt := time.Now()
	user.UpdatedAt = newUpdatedAt
	user.CreatedAt = newCreatedAt

	// Test updated values
	assert.WithinDuration(t, newUpdatedAt, user.GetUpdatedAt(), time.Second)
	assert.WithinDuration(t, newCreatedAt, user.GetCreatedAt(), time.Second)
}
