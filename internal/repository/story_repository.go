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

	if f.Size != "" {
		query += " AND size = ?"
		args = append(args, f.Size)
	}
	if f.AIGenerated != nil {
		query += " AND ai_generated = ?"
		if *f.AIGenerated {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
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
		var aiGen int
		if err := rows.Scan(
			&s.ID, &s.Title, &s.CoverImage, &s.Author, &s.Content,
			&aiGen, &s.Size, &s.Views, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		s.AIGenerated = aiGen == 1
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

	if _, err := tx.ExecContext(ctx, `UPDATE stories SET views = views + 1 WHERE id = ?`, id); err != nil {
		return nil, err
	}

	var s models.Story
	var aiGen int
	err = tx.QueryRowContext(ctx, `
		SELECT id, title, cover_image, author, content, ai_generated, size, views, created_at, updated_at
		FROM stories WHERE id = ?`, id).Scan(
		&s.ID, &s.Title, &s.CoverImage, &s.Author, &s.Content,
		&aiGen, &s.Size, &s.Views, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	s.AIGenerated = aiGen == 1
	return &s, tx.Commit()
}

func (r *StoryRepository) Create(ctx context.Context, s *models.Story) (*models.Story, error) {
	now := time.Now().UTC()
	aiGen := boolToInt(s.AIGenerated)

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO stories (title, cover_image, author, content, ai_generated, size, views, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?)`,
		s.Title, s.CoverImage, s.Author, s.Content, aiGen, s.Size, now, now,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.getByIDRaw(ctx, id)
}

func (r *StoryRepository) Update(ctx context.Context, id int64, s *models.Story) (*models.Story, error) {
	now := time.Now().UTC()
	aiGen := boolToInt(s.AIGenerated)

	res, err := r.db.ExecContext(ctx, `
		UPDATE stories
		SET title=?, cover_image=?, author=?, content=?, ai_generated=?, size=?, updated_at=?
		WHERE id=?`,
		s.Title, s.CoverImage, s.Author, s.Content, aiGen, s.Size, now, id,
	)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, nil
	}
	return r.getByIDRaw(ctx, id)
}

func (r *StoryRepository) Delete(ctx context.Context, id int64) (bool, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM stories WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// getByIDRaw fetches a story without touching the view counter.
func (r *StoryRepository) getByIDRaw(ctx context.Context, id int64) (*models.Story, error) {
	var s models.Story
	var aiGen int
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, cover_image, author, content, ai_generated, size, views, created_at, updated_at
		FROM stories WHERE id = ?`, id).Scan(
		&s.ID, &s.Title, &s.CoverImage, &s.Author, &s.Content,
		&aiGen, &s.Size, &s.Views, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("getByIDRaw: %w", err)
	}
	s.AIGenerated = aiGen == 1
	return &s, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
