package dbconn

import (
	. "context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wujiu2020/strip/utils"
	"github.com/wujiu2020/strip/utils/helpers"
)

type Connect struct {
	name  ConnectorName
	getDB dbGetter
}

func (p *Connect) Invoke(ctx Context, handle interface{}) error {

	var (
		isDB bool
		isTx bool
		isQ  bool
		hDB  handleDB
		hTX  handleTX
		hQ   handleQ
	)
	switch v := handle.(type) {
	case func(*sqlx.DB) error:
		isDB = true
		hDB = v
	case func(*sqlx.Tx) error:
		isTx = true
		hTX = v
	case func(Q) error:
		isQ = true
		hQ = v
	}

	var db *sqlx.DB
	var tx *sqlx.Tx
	utils.CtxFindValue(ctx, &db)
	utils.CtxFindValue(ctx, &tx)

	switch {
	case isDB:
		if db == nil {
			db = p.getDB()
		}
	case isTx:
		if tx != nil {
			return hTX(tx)
		}
	case isQ:
		if tx != nil {
			return hQ(tx)
		}
		if db == nil {
			db = p.getDB()
		}
	}

	switch {
	case isDB:
		return hDB(db)
	case isTx:
		return p.invokeTx(ctx, hTX)
	case isQ:
		return hQ(db)
	default:
		return fmt.Errorf("not support db handler %T", handle)
	}
}

func (p *Connect) invokeTx(ctx Context, handle handleTX) (err error) {
	// Begin a transcation with context.
	tx, err := p.getDB().BeginTxx(ctx, nil)
	if err != nil {
		return
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
			er := tx.Rollback()
			helpers.GetLogger(ctx).Warn("rollback:", er)
			return
		}

		// Commit
		err = tx.Commit()
		return
	}()

	err = handle(tx)
	return
}

func WithDB(ctx Context, db interface{}) Context {
	switch db := db.(type) {
	case *sqlx.DB:
		ctx = utils.CtxWithValue(ctx, db)
	case *sqlx.Tx:
		ctx = utils.CtxWithValue(ctx, db)
	default:
		panic(fmt.Sprintf("not support db %T", db))
	}
	return ctx
}
