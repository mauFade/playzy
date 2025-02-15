package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/dto"
	"github.com/mauFade/playzy/internal/model"
)

type SessionRepositoryInterface interface {
	Create(s *model.SessionModel) error
	FindByID(id uuid.UUID) (*model.SessionModel, error)
	FindAvailable(page int, rank, game string) (*dto.SessionsPageResponse, error)
	Delete(id string) error
}

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(d *sql.DB) *SessionRepository {
	r := &SessionRepository{
		db: d,
	}

	r.db.Exec("CREATE TABLE IF NOT EXISTS sessions (id UUID PRIMARY KEY, game VARCHAR NOT NULL, user_id UUID NOT NULL, objective VARCHAR NOT NULL, rank VARCHAR NULL, is_ranked BOOLEAN NOT NULL, updated_at TIMESTAMP NOT NULL, created_at TIMESTAMP NOT NULL, CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE)")

	return r
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

func (r *SessionRepository) FindAvailable(page int, rank, game string) (*dto.SessionsPageResponse, error) {
	pageQtd := 6

	rankQuery := ""
	gameQuery := ""

	if rank != "" {
		rankQuery = "AND LOWER(sessions.rank) = '" + strings.ToLower(rank) + "'"
	}

	if game != "" {
		gameQuery = "AND LOWER(sessions.game) LIKE '%" + strings.ToLower(game) + "%'"
	}

	query := fmt.Sprintf("SELECT sessions.*, users.id AS user_id, users.name, users.email, users.gamertag FROM sessions LEFT JOIN users ON sessions.user_id = users.id WHERE users.is_deleted = 'false' %s %s LIMIT %v OFFSET ($1 - 1) * %v",
		rankQuery, gameQuery, pageQtd, pageQtd)

	fmt.Println("\n\n", query)
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

func (r *SessionRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id = $1", id)

	return err
}
