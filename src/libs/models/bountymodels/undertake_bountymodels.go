package bountymodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/models/ethmodels"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
	ethSdk "cos-backend-com/src/libs/sdk/eth"
	"github.com/jmoiron/sqlx"
	"time"
)

var UndertakeBounties = &undertakeBounties{
	Connector: models.DefaultConnector,
}

type undertakeBounties struct {
	dbconn.Connector
}

func (c *undertakeBounties) CreateUndertakeBounty(ctx context.Context, uid flake.ID, input *coresSdk.CreateUndertakeBountyInput, output *coresSdk.UndertakeBountyResult) (err error) {
	stmt := `
		INSERT INTO bounties_hunters_rel(bounty_id, uid, status, started_at)
		VALUES (${bountyId}, ${uid}, ${status}, ${startedAt}) RETURNING id, bounty_id, status;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{bountyId}":  input.BountyId,
		"{uid}":       uid,
		"{status}":    coresSdk.UndertakeBountyStatusNull,
		"{startedAt}": time.Now(),
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		if er := db.GetContext(newCtx, output, query, args...); er != nil {
			return er
		}
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceUndertakeBounty,
			SourceId: output.Id,
		}

		return ethmodels.Transactions.Create(newCtx, &createTransactionsInput)
	})
}
