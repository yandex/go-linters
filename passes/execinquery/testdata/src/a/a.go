package a

import (
	"context"
	"database/sql"
	"errors"
)

const selectWithComment = `-- foobar
SELECT * FROM test WHERE test=?
`

const deleteWithComment = `-- foobar
-- foobar
-- foobar
DELETE * FROM test WHERE test=?
`

const deleteWithCommentMultiline = `/* foobar
-- foobar
-- foobar
   */
DELETE * FROM test WHERE test=?
`

func sample(db *sql.DB) {
	s := "alice"

	_ = db.QueryRowContext(context.Background(), "SELECT * FROM test WHERE test=?", s)
	_ = db.QueryRowContext(context.Background(), selectWithComment, s)
	_ = db.QueryRowContext(context.Background(), deleteWithComment, s)          // want "Use ExecContext instead of QueryRowContext to execute `DELETE` query"
	_ = db.QueryRowContext(context.Background(), deleteWithCommentMultiline, s) // want "Use ExecContext instead of QueryRowContext to execute `DELETE` query"

	_ = db.QueryRowContext(context.Background(), "DELETE * FROM test WHERE test=?", s) // want "Use ExecContext instead of QueryRowContext to execute `DELETE` query"
	_ = db.QueryRowContext(context.Background(), "UPDATE * FROM test WHERE test=?", s) // want "Use ExecContext instead of QueryRowContext to execute `UPDATE` query"

	_, _ = db.Query("UPDATE * FROM test WHERE test=?", s)                              // want "Use Exec instead of Query to execute `UPDATE` query"
	_, _ = db.QueryContext(context.Background(), "UPDATE * FROM test WHERE test=?", s) // want "Use ExecContext instead of QueryContext to execute `UPDATE` query"
	_ = db.QueryRow("UPDATE * FROM test WHERE test=?", s)                              // want "Use Exec instead of QueryRow to execute `UPDATE` query"

	_, _ = db.Query(otherFileValue, s)

	query := "UPDATE * FROM test where test=?"
	_, _ = db.Query(query, s) // want "Use Exec instead of Query to execute `UPDATE` query"

	f1 := `
UPDATE * FROM test WHERE test=?`
	_ = db.QueryRow(f1, s) // want "Use Exec instead of QueryRow to execute `UPDATE` query"

	const f2 = `
UPDATE * FROM test WHERE test=?`
	_ = db.QueryRow(f2, s) // want "Use Exec instead of QueryRow to execute `UPDATE` query"

	f3 := `
UPDATE * FROM test WHERE test=?`
	_ = db.QueryRow(f3, s) // want "Use Exec instead of QueryRow to execute `UPDATE` query"

	f4 := f3
	_ = db.QueryRow(f4, s) // want "Use Exec instead of QueryRow to execute `UPDATE` query"

	f5 := `
UPDATE * ` + `FROM test` + ` WHERE test=?`
	_ = db.QueryRow(f5, s) // want "Use Exec instead of QueryRow to execute `UPDATE` query"

	err := errors.New("oops")
	err.Error()
}
