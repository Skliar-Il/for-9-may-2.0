package repository

import (
	"context"
	"for9may/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GalleryRepositoryInterface interface {
	Create(ctx context.Context, tx pgx.Tx, post *dto.CreateGalleryPostDTO) (*uuid.UUID, error)
	Delete(ctx context.Context, tx pgx.Tx, postID uuid.UUID) error
	GetAll(ctx context.Context, tx pgx.Tx) ([]*dto.GalleryPostDTO, error)
	UpdateLink(ctx context.Context, tx pgx.Tx, postID uuid.UUID, link string) error
}

type GalleryRepository struct{}

func NewGalleryRepository() *GalleryRepository {
	return &GalleryRepository{}
}

func (GalleryRepository) Create(ctx context.Context, tx pgx.Tx, post *dto.CreateGalleryPostDTO) (*uuid.UUID, error) {
	query := `
		INSERT INTO gallery(date, description)
		VALUES($1, $2)
		RETURNING id
	`

	var id uuid.UUID
	if err := tx.QueryRow(ctx, query, post.Date, post.Description).Scan(&id); err != nil {
		return nil, err
	}

	return &id, nil
}

func (GalleryRepository) Delete(ctx context.Context, tx pgx.Tx, postID uuid.UUID) error {
	query := `
	DELETE FROM gallery
	WHERE id = $1
	`

	_, err := tx.Exec(ctx, query, postID)
	return err
}

func (GalleryRepository) GetAll(ctx context.Context, tx pgx.Tx) ([]*dto.GalleryPostDTO, error) {
	query := `
		SELECT id, description, COALESCE(link, ''), date
		FROM gallery
	`

	var posts []*dto.GalleryPostDTO
	postRows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for postRows.Next() {
		var post dto.GalleryPostDTO
		if err := postRows.Scan(&post.ID, &post.Description, &post.Link, &post.Date); err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (GalleryRepository) UpdateLink(ctx context.Context, tx pgx.Tx, postID uuid.UUID, link string) error {
	query := `
	UPDATE gallery
	SET 
		link = $1
	WHERE id = $2
	`

	_, err := tx.Exec(ctx, query, link, postID)
	return err
}
