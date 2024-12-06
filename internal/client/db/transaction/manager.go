package transaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/sSmok/auth/internal/client/db"
	"github.com/sSmok/auth/internal/client/db/pg"
)

type manager struct {
	pool db.TransactorI
}

func NewManager(pool db.TransactorI) db.TxManagerI {
	return &manager{pool: pool}
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	opts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, opts, f)
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, f db.Handler) (err error) {
	// Если это вложенная транзакция, пропускаем инициацию новой транзакции и выполняем обработчик.
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return f(ctx)
	}

	// Стартуем новую транзакцию.
	tx, err = m.pool.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	// Кладем транзакцию в контекст.
	ctx = pg.MakeContextTransaction(ctx, tx)

	// Настраиваем функцию отсрочки для отката или коммита транзакции.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		// откатываем транзакцию, если произошла ошибка
		if err != nil {
			if errRallback := tx.Rollback(ctx); errRallback != nil {
				err = errors.Wrap(errRallback, "rollback transaction")
			}
			return
		}

		// если ошибок не было, коммитим транзакцию
		err = tx.Commit(ctx)
		if err != nil {
			err = errors.Wrap(err, "transaction commit failed")
		}

	}()

	// Выполните код внутри транзакции.
	// Если функция терпит неудачу, возвращаем ошибку, и функция отсрочки выполняет откат
	// или в противном случае транзакция коммитится.
	if err = f(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return
}
