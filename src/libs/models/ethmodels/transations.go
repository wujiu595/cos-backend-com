package ethmodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	ethSdk "cos-backend-com/src/libs/sdk/eth"

	"github.com/jmoiron/sqlx"
)

var Transactions = &transactions{
	Connector: models.DefaultConnector,
}

type transactions struct {
	dbconn.Connector
}

func (c *transactions) Insert(ctx context.Context, input *ethSdk.CreateTransactionsInput) (err error) {
	stmt := `
		INSERT INTO transactions(tx_id, source, source_id)
		VALUES (${txId}, ${source}, ${sourceId})
		ON CONFLICT(tx_id) WHERE state !=3 
		DO UPDATE SET 
		(
		    state, retry_time, source_id, updated_at
		)= (
		    1, 0, ${sourceId}, current_timestamp
		);
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{source}":   input.Source,
		"{sourceId}": input.SourceId,
		"{txId}":     input.TxId,
	})

	return c.Invoke(ctx, func(db *sqlx.DB) (er error) {
		_, err = db.ExecContext(ctx, query, args...)
		return err
	})
}
