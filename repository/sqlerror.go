package db

import "github.com/jackc/pgx/v5/pgconn"

var (
	pgErr *pgconn.PgError
)

type SqlDuplicateEntryError struct {
	msg string
}
func (e SqlDuplicateEntryError) Error() string {
	return e.msg
}

func newSqlDuplicateEntryError(msg string) error {
	return &SqlDuplicateEntryError{
		msg: msg,
	}
}
