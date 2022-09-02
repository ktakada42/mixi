package testutil

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

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
