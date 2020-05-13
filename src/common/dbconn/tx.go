package dbconn

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NestTx(ctx context.Context, tx *sqlx.Tx, handle handleTX) (err error) {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	savepoint := string(b)
	_, err = tx.ExecContext(ctx, "SAVEPOINT "+savepoint)
	if err != nil {
		return err
	}
	defer func() {
		if re := recover(); re != nil {
			switch re := re.(type) {
			case error:
				err = re
			default:
				err = fmt.Errorf("%s", re)
			}
		}

		// Check Error & Rollback
		if err != nil {
			_, xerr := tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT "+savepoint)
			if xerr != nil {
				err = fmt.Errorf("rollback to savepoint error: %v while meeting %v", xerr, err)
			}
			return
		}

		_, xerr := tx.ExecContext(ctx, "RELEASE SAVEPOINT "+savepoint)
		if xerr != nil {
			err = fmt.Errorf("release savepoint error: %v", xerr)
		}
		return
	}()

	err = handle(tx)
	return
}
