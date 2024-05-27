package data

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lib/pq"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Seed: Seed{},
	}
}

type Models struct {
	Seed Seed
}

type Seed struct {
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *Seed) GetAll(limit int) ([]*Seed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT 
		url,
		processed,
		created_at,
		updated_at
	FROM seeds ORDER BY created_at asc LIMIT $1`

	rows, err := db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seeds []*Seed

	for rows.Next() {
		var seed Seed
		err := rows.Scan(
			&seed.Url,
			&seed.CreatedAt,
			&seed.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		seeds = append(seeds, &seed)
	}

	return seeds, nil
}

func (s *Seed) SaveProcessed(seeds []*Seed) error {
	var urls []string
	for _, seed := range seeds {
		urls = append(urls, seed.Url)
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE seeds SET processed = true WHERE url = ANY($1)`
	_, err := db.ExecContext(ctx, query, pq.Array(urls))
	if err != nil {
		log.Println("Error updating", err)
		return err
	}

	return nil
}
