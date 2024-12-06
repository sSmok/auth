package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ClientI интерфейс для работы с БД
type ClientI interface {
	DB() DB
	Close() error
}

// Query структура для выполнения запросов, хранящая имя запроса и сам запрос
// Имя запроса используется для логирования и потенциально может использоваться еще где-то, например, для трейсинга
type Query struct {
	Name     string
	QueryRaw string
}

// Handler - функция, которая выполняется в транзакции
type Handler func(ctx context.Context) error

// TxManagerI интерфейс выполнения переданного обработчика в рамках определенного уровня изоляции
type TxManagerI interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// TransactorI интерфейс для работы с транзакцией
type TransactorI interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// DB интерфейс для запросов к БД
type DB interface {
	Pinger
	SQLExecer
	TransactorI
	Close()
}

// SQLExecer интрефейс для выполнения запросов
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer интрефейс для выполнения именованных запросов с помощью тегов в структурах
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, query Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, query Query, args ...interface{}) error
}

// QueryExecer интрефейс для выполнения запросов
type QueryExecer interface {
	ExecContext(ctx context.Context, query Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, query Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, query Query, args ...interface{}) pgx.Row
}

// Pinger интерфейс для проверки подключения
type Pinger interface {
	Ping(ctx context.Context) error
}
