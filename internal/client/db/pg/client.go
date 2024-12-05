package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/auth/internal/client/db"
)

type pgClient struct {
	masterDB db.DB
}

// NewPGClient создает новый клиент для работы с БД через переданный DSN
func NewPGClient(ctx context.Context, dsn string) (db.ClientI, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &pgClient{
		masterDB: NewPG(pool),
	}, nil
}

func (p *pgClient) DB() db.DB {
	return p.masterDB
}

func (p *pgClient) Close() error {
	if p.masterDB != nil {
		p.masterDB.Close()
	}

	return nil
}
