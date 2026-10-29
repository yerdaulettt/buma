package postgres

import (
	"context"

	"buma/internal/image"

	"github.com/jackc/pgx/v5/pgxpool"
)

type imageRepo struct {
	db *pgxpool.Pool
}

func NewImageRepo(db *pgxpool.Pool) *imageRepo {
	return &imageRepo{db: db}
}

func (r *imageRepo) AddMetadata(ctx context.Context, userId int, imageType, key string) error {
	query := `insert into images (image_type, s3_key, user_id) values ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, imageType, key, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *imageRepo) GetMetadata(ctx context.Context, userId int) ([]image.ImageMetadata, error) {
	query := `select id, added_at, s3_key, user_id from images where user_id = $1`

	var mdata []image.ImageMetadata

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m image.ImageMetadata

		if err := rows.Scan(&m.Id, &m.AddedAt, &m.Key, &m.UserId); err != nil {
			return nil, err
		}

		mdata = append(mdata, m)
	}

	return mdata, nil
}
