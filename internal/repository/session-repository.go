package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/dto"
	"github.com/mauFade/playzy/internal/model"
)

type SessionRepositoryInterface interface {
	Create(s *model.SessionModel) error
	FindByID(id uuid.UUID) (*model.SessionModel, error)
	FindAvailable(page int) (*dto.SessionsPageResponse, error)
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

func (r *SessionRepository) FindAvailable(page int) (*dto.SessionsPageResponse, error) {
	pageQtd := 4

	query := fmt.Sprintf("SELECT sessions.*, users.id AS user_id, users.name, users.email, users.gamertag FROM sessions LEFT JOIN users ON sessions.user_id = users.id WHERE users.is_deleted = 'false' LIMIT %v OFFSET ($1 - 1) * %v",
		pageQtd, pageQtd)

	fmt.Println(query)
	rows, err := r.db.Query(query, page)

	if err != nil {
		return nil, err
	}

	var count int

	err = r.db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&count)

	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(pageQtd)))

	var sessions []dto.SessionWithUser

	for rows.Next() {
		var s dto.SessionWithUser

		err := rows.Scan(&s.ID, &s.Game, &s.UserID, &s.Objective, &s.Rank, &s.IsRanked, &s.UpdatedAt, &s.CreatedAt, &s.UserID, &s.UserName, &s.Email, &s.UserGamertag)

		if err != nil {
			return nil, err
		}

		sessions = append(sessions, s)
	}

	if len(sessions) == 0 {
		sessions = []dto.SessionWithUser{}
	}

	return &dto.SessionsPageResponse{
		Page:       page,
		TotalPages: totalPages,
		Sessions:   sessions,
	}, nil
}
