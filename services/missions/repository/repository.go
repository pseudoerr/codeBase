package repository

import (
	"context"
	"database/sql"
	"github.com/pseudoerr/mission-service/models"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) ListMissions(ctx context.Context) ([]models.Mission, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, title, points FROM missions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var m models.Mission
		if err := rows.Scan(&m.ID, &m.Title, &m.Points); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	return missions, nil
}

func (r *PostgresRepository) AddMission(ctx context.Context, m models.Mission) (models.Mission, error) {
	err := r.DB.QueryRowContext(
		ctx,
		"INSERT INTO missions (title, points) VALUES ($1, $2) RETURNING id",
		m.Title, m.Points,
	).Scan(&m.ID)

	return m, err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (models.Mission, error) {
	var m models.Mission
	err := r.DB.QueryRowContext(ctx, "SELECT id, title, points FROM missions WHERE id = $1", id).
		Scan(&m.ID, &m.Title, &m.Points)
	return m, err
}

func (r *PostgresRepository) UpdateMission(ctx context.Context, m models.Mission) (models.Mission, error) {
	_, err := r.DB.ExecContext(
		ctx,
		"UPDATE missions SET title = $1, points = $2 WHERE id = $3",
		m.Title, m.Points, m.ID,
	)
	return m, err
}

func (r *PostgresRepository) DeleteMission(ctx context.Context, id int) error {
	_, err := r.DB.Exec("DELETE FROM missions WHERE id = $1", id)
	return err
}
