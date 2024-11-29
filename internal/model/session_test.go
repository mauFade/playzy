package model_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewSessionModel(t *testing.T) {
	id := uuid.New()
	game := "VALORANT"
	userID := uuid.New()
	objective := "Serious play"
	rank := "Diamont"
	isRanked := true
	updatedAt := time.Now()
	createdAt := time.Now()

	session := model.NewSessionModel(
		id,
		userID,
		game,
		objective,
		&rank,
		isRanked,
		updatedAt,
		createdAt,
	)

	assert.Equal(t, id, session.GetID())
	assert.Equal(t, game, session.GetGame())
	assert.Equal(t, userID, session.GetUserID())
	assert.Equal(t, objective, session.GetObjective())
	assert.Equal(t, &rank, session.GetRank())
	assert.Equal(t, isRanked, session.GetIsRanked())
	assert.WithinDuration(t, updatedAt, session.GetUpdatedAt(), time.Second)
	assert.WithinDuration(t, createdAt, session.GetCreatedAt(), time.Second)
}

func TestSessionSettersAndGetters(t *testing.T) {
	session := model.SessionModel{}
	session.IsRanked = true
	session.Rank = nil

	userID := uuid.New()
	rank := "Diamont"

	session.SetUserID(userID)
	session.SetGame("VALORANT")
	session.SetObjective("Serious play")
	session.SetRank(&rank)
	session.SetIsRankedOrNot()

	assert.Equal(t, userID, session.GetUserID())
	assert.Equal(t, "VALORANT", session.GetGame())
	assert.Equal(t, "Serious play", session.GetObjective())
	assert.Equal(t, &rank, session.GetRank())
	assert.Equal(t, false, session.GetIsRanked())
}
