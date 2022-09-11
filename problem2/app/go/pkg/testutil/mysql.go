package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql" // mysql

	"github.com/DATA-DOG/go-txdb"
)

var txDBRegisterOnce sync.Once

func registerTxDB() {
	txdb.Register("txdb", "mysql", "root:@(localhost:3306)/app")
}

func PrepareMySQL(t *testing.T) *sql.DB {
	t.Helper()
	txDBRegisterOnce.Do(registerTxDB)

	cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
	db, err := sql.Open("txdb", cName)
	if err != nil {
		t.Fatalf("failed to open txdb connection: %s", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func ValidateSQLArgs(t *testing.T, q string, args ...any) {
	t.Helper()

	numQ := strings.Count(q, "?")
	if numQ != len(args) {
		t.Fatalf("invalid args: q=%s, numQ=%d, args=%v, len(args)=%d", q, numQ, args, len(args))
	}
}

func ExecSQL(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()

	if _, err := db.Exec(q, args...); err != nil {
		t.Fatal(err)
	}
}

func NewSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	return db, mock
}

func BeginTx(t *testing.T, db *sql.DB) *sql.Tx {
	t.Helper()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		t.Fatal(err)
	}

	return tx
}

func CommitTx(t *testing.T, tx *sql.Tx) {
	t.Helper()

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func RollBackTx(t *testing.T, tx *sql.Tx) {
	t.Helper()

	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}
