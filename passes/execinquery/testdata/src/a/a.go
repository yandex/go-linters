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

const showWithComment = `-- foobar
SHOW tables
`

func sample(db *sql.DB) {
	s := "alice"

	_ = db.QueryRowContext(context.Background(), "SELECT * FROM test WHERE test=?", s)
	_ = db.QueryRowContext(context.Background(), selectWithComment, s)

	// SHOW queries should be allowed (they return results)
	_, _ = db.Query("SHOW tables")
	_, _ = db.Query("SHOW databases")
	_, _ = db.QueryContext(context.Background(), "SHOW tables")
	_ = db.QueryRow("SHOW tables")
	_ = db.QueryRowContext(context.Background(), "SHOW tables")
	_ = db.QueryRowContext(context.Background(), showWithComment)
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

	// PostgreSQL RETURNING queries should be allowed (they return results)
	// INSERT with RETURNING
	_ = db.QueryRow("INSERT INTO test (name) VALUES (?) RETURNING id", s)
	_ = db.QueryRowContext(context.Background(), "INSERT INTO test (name) VALUES (?) RETURNING id", s)
	_, _ = db.Query("INSERT INTO test (name) VALUES (?) RETURNING id, name", s)
	_, _ = db.QueryContext(context.Background(), "INSERT INTO test (name) VALUES (?) RETURNING *", s)

	// UPDATE with RETURNING
	_ = db.QueryRow("UPDATE test SET name=? WHERE id=1 RETURNING id", s)
	_ = db.QueryRowContext(context.Background(), "UPDATE test SET name=? WHERE id=1 RETURNING id, name", s)
	_, _ = db.Query("UPDATE test SET name=? WHERE id=1 RETURNING *", s)
	_, _ = db.QueryContext(context.Background(), "UPDATE test SET name=? RETURNING id", s)

	// DELETE with RETURNING
	_ = db.QueryRow("DELETE FROM test WHERE id=1 RETURNING id", s)
	_ = db.QueryRowContext(context.Background(), "DELETE FROM test WHERE id=1 RETURNING id, name", s)
	_, _ = db.Query("DELETE FROM test WHERE id=1 RETURNING *", s)
	_, _ = db.QueryContext(context.Background(), "DELETE FROM test WHERE id=1 RETURNING id", s)

	// RETURNING with different casing and formatting
	_ = db.QueryRow("INSERT INTO test (name) VALUES (?) returning id", s)
	_ = db.QueryRow("UPDATE test SET name=? WHERE id=1 ReTuRnInG id", s)
	_ = db.QueryRow(`
		INSERT INTO test (name)
		VALUES (?)
		RETURNING id, created_at
	`, s)

	// RETURNING in const/var
	const insertReturning = "INSERT INTO test (name) VALUES (?) RETURNING id"
	_ = db.QueryRow(insertReturning, s)

	updateReturning := "UPDATE test SET name=? WHERE id=1 RETURNING id"
	_ = db.QueryRow(updateReturning, s)

	// RETURNING with comments
	const insertReturningWithComment = `-- Insert and return ID
	INSERT INTO test (name) VALUES (?) RETURNING id`
	_ = db.QueryRow(insertReturningWithComment, s)

	// UPSERT with RETURNING (ON CONFLICT ... DO UPDATE)
	upsertQuery := `
		INSERT INTO users (id, other_id, display_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (other_id)
		DO UPDATE SET display_name = COALESCE(NULLIF(EXCLUDED.display_name, ''), users.display_name)
		RETURNING id, other_id, display_name;
	`
	_ = db.QueryRow(upsertQuery, 1, 2, "test")

	err := errors.New("oops")
	err.Error()
}
