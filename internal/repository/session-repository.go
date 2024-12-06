package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
)

type SessionRepositoryInterface interface {
	Create(s *model.SessionModel) error
	FindByID(id uuid.UUID) (*model.SessionModel, error)
}

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(d *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: d,
	}
}

func (r *SessionRepository) Create(s *model.SessionModel) error {
	query := `INSERT INTO sessions
	(id, game, user_id, objective, rank, is_ranked, updated_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := r.db.Exec(query,
		s.GetID(),
		s.GetGame(),
		s.GetUserID(),
		s.GetObjective(),
		s.GetRank(),
		s.GetIsRanked(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) FindByID(id uuid.UUID) (*model.SessionModel, error) {
	query := `SELECT * FROM sessions WHERE id = $1`

	res := r.db.QueryRow(query, id.String())

	var session model.SessionModel

	if err := res.Scan(
		&session.ID,
		&session.Game,
		&session.UserID,
		&session.Objective,
		&session.Rank,
		&session.IsRanked,
		&session.UpdatedAt,
		&session.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &session, nil
}
