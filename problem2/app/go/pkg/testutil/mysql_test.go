package testutil

import "testing"

func Test_mysql_TxActions(t *testing.T) {
	db := PrepareMySQL(t)
	defer db.Close()

	tx := BeginTx(t, db)
	CommitTx(t, tx)

	tx2 := BeginTx(t, db)
	RollBackTx(t, tx2)
}
