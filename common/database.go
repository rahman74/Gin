package common

import (
	"database/sql"
)

// TxChecker ...
func TxChecker(tx *sql.Tx, err error) {
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}
}
