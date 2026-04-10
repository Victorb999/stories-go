package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"stories-go/internal/models"
)

type StoryRepository struct {
	db *sql.DB
}

func NewStoryRepository(db *sql.DB) *StoryRepository {
	return &StoryRepository{db: db}
}

// ListFilter holds optional filters for the List query.
type ListFilter struct {
	Size        string
	AIGenerated *bool
}

func (r *StoryRepository) List(ctx context.Context, f ListFilter) ([]models.Story, error) {
	query := `SELECT id, title, cover_image, author, content, ai_generated, size, views, created_at, updated_at
	          FROM stories WHERE 1=1`
	args := []any{}
	paramID := 1

	if f.Size != "" {
		query += fmt.Sprintf(" AND size = $%d", paramID)
		args = append(args, f.Size)
		paramID++
	}
	if f.AIGenerated != nil {
		query += fmt.Sprintf(" AND ai_generated = $%d", paramID)
		args = append(args, *f.AIGenerated)
		paramID++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []models.Story
	for rows.Next() {
		var s models.Story
		if err := rows.Scan(
			&s.ID, &s.Title, &s.CoverImage, &s.Author, &s.Content,
			&s.AIGenerated, &s.Size, &s.Views, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		stories = append(stories, s)
	}
	return stories, rows.Err()
}

// GetByID fetches a story and atomically increments its view counter.
func (r *StoryRepository) GetByID(ctx context.Context, id int64) (*models.Story, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `UPDATE stories SET views = views + 1 WHERE id = $1`, id); err != nil {
		return nil, err
	}

	var s models.Story
	err = tx.QueryRowContext(ctx, `
		SELECT id, title, cover_image, author, content, ai_generated, size, views, created_at, updated_at
		FROM stories WHERE id = $1`, id).Scan(
		&s.ID, &s.Title, &s.CoverImage, &s.Author, &s.Content,
		&s.AIGenerated, &s.Size, &s.Views, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, tx.Commit()
}

func (r *StoryRepository) Create(ctx context.Context, s *models.Story) (*models.Story, error) {
	now := time.Now().UTC()

	var newID int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO stories (title, cover_image, author, content, ai_generated, size, views, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, 0, $7, $8) RETURNING id`,
		s.Title, s.CoverImage, s.Author, s.Content, s.AIGenerated, s.Size, now, now,
	).Scan(&newID)

	if err != nil {
		return nil, err
	}
	return r.getByIDRaw(ctx, newID)
}

func (r *StoryRepository) Update(ctx context.Context, id int64, s *models.Story) (*models.Story, error) {
	now := time.Now().UTC()

	res, err := r.db.ExecContext(ctx, `
		UPDATE stories
		SET title=$1, cover_image=$2, author=$3, content=$4, ai_generated=$5, size=$6, updated_at=$7
		WHERE id=$8`,
		s.Title, s.CoverImage, s.Author, s.Content, s.AIGenerated, s.Size, now, id,
	)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, nil // Not found
	}
	return r.getByIDRaw(ctx, id)
}

func (r *StoryRepository) Delete(ctx context.Context, id int64) (bool, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM stories WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// getByIDRaw fetches a story without touching the view counter.
func (r *StoryRepository) getByIDRaw(ctx context.Context, id int64) (*models.Story, error) {
	var s models.Story
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, cover_image, author, content, ai_generated, size, views, created_at, updated_at
		FROM stories WHERE id = $1`, id).Scan(
		&s.ID, &s.Title, &s.CoverImage, &s.Author, &s.Content,
		&s.AIGenerated, &s.Size, &s.Views, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("getByIDRaw: %w", err)
	}
	return &s, nil
}
