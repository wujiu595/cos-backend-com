package bountymodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/dbquery"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/types"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/models/ethmodels"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
	ethSdk "cos-backend-com/src/libs/sdk/eth"
	"time"

	"github.com/jmoiron/sqlx"
)

var Bounties = &bounties{
	Connector: models.DefaultConnector,
}

// categories represents controller for 'categories'.
type bounties struct {
	dbconn.Connector
}

func (c *bounties) CreateBounty(ctx context.Context, startupId, uid flake.ID, input *coresSdk.CreateBountyInput) (err error) {
	stmt := `
		INSERT INTO bounties(id, startup_id, user_id, title, type, keywords, contact_email, intro, description_addr, description_file_addr,
			duration, expired_at, payments)
		VALUES (${id}, ${startupId}, ${userId}, ${title}, ${type}, ARRAY [${keywords}], ${contactEmail}, ${intro}, ${descriptionAddr},
			${descriptionFileAddr}, ${duration}, ${expiredAt}, ${payments});
	`

	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{id}":                  input.Id,
		"{startupId}":           startupId,
		"{userId}":              uid,
		"{title}":               input.Title,
		"{type}":                input.Type,
		"{keywords}":            input.Keywords,
		"{contactEmail}":        input.ContactEmail,
		"{intro}":               input.Intro,
		"{descriptionAddr}":     input.DescriptionAddr,
		"{descriptionFileAddr}": input.DescriptionFileAddr,
		"{duration}":            input.Duration,
		"{expiredAt}":           time.Now().AddDate(0, 0, input.Duration),
		"{payments}":            types.JSONAny{input.Payments},
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceBounty,
			SourceId: input.Id,
		}

		if err := ethmodels.Transactions.Create(newCtx, &createTransactionsInput); err != nil {
			return err
		}
		_, err = db.ExecContext(newCtx, query, args...)
		return err
	})
}

// List Categories interface{} by input
func (c *bounties) ListBounties(ctx context.Context, startupId, uid flake.ID, isOwner bool, input *coresSdk.ListBountiesInput, outputs interface{}) (total int, err error) {
	plan := &dbquery.Plan{}
	plan.RetTotal = true

	filterStmt := ""
	keyword := ""
	if input.Keyword != "" {
		keyword = "%" + util.PgEscapeLike(input.Keyword) + "%"
		filterStmt += `AND b.name ILIKE ${keyword}`
	}

	if startupId != flake.ID(0) {
		plan.AddCond(`AND b.startup_id = ${startupId}`)
	}

	plan.OrderBySql = ` ORDER BY is_open ASC,created_at DESC`
	plan.LimitSql = ` LIMIT ${limit} OFFSET ${offset}`

	plan.Params = map[string]interface{}{
		"{keyword}":   keyword,
		"{uid}":       uid,
		"{startupId}": startupId,
		"{offset}":    input.Offset,
		"{limit}":     input.GetLimit(),
	}

	total, err = c.Query(ctx, uid, isOwner, outputs, plan)
	return
}

func (c *bounties) Query(ctx context.Context, uid flake.ID, isOwner bool, m interface{}, plan *dbquery.Plan) (total int, err error) {
	filterSql := `
        FROM bounties b
        INNER JOIN startups s ON s.id = b.startup_id
	`
	joinCondition := ``

	if uid != 0 {
		joinCondition += "INNER JOIN bounties_hunters_rel bhr ON bhr.bounty_id = b.id AND bhr.uid = ${uid}"
	}

	if isOwner {
		joinCondition += "LEFT JOIN transactions t ON b.id = t.source_id AND t.source = ${source}"
	} else {
		joinCondition += "INNER JOIN transactions t ON b.id = t.source_id AND t.source = ${source}"
	}

	plan.Params["{source}"] = ethSdk.TransactionSourceBounty

	if plan.RetTotal {
		query :=
			`SELECT COUNT(*)
			` + filterSql + `
			` + joinCondition + `
			WHERE 1=1
			` + plan.Conditions
		query, args := util.PgMapQueryV2(query, plan.Params)

		err = c.Invoke(ctx, func(db dbconn.Q) error {
			return db.GetContext(ctx, &total, query, args...)
		})
		if err != nil {
			return
		}
	}

	query := `
	WITH bounties_cte AS (
		SELECT b.*,t.block_addr, t.state transaction_state, current_timestamp<b.expired_at is_open,json_build_object('id',s.id,'name',s.name) startup
		` + filterSql + `
        ` + joinCondition + `
		WHERE 1=1 
		` + plan.Conditions + `
		` + plan.OrderBySql + `
		` + plan.LimitSql + `
	),bounty_hunter_rels_cte AS (
		SELECT bhr.bounty_id,h.id hunter_id, bhr.uid as user_id, bhr.status, bhr.started_at, bhr.submitted_at, bhr.quited_at, bhr.paid_at, bhr.paid_tokens,COALESCE(h.name, u.public_key) AS name
		FROM bounties_cte bc
		LEFT JOIN bounties_hunters_rel bhr ON bhr.bounty_id = bc.id
		LEFT JOIN users u ON bhr.uid = u.id
		LEFT JOIN hunters h ON u.id = h.user_id
	),bounty_hunter_rels_aggregate_cte AS (
		SELECT bhrc.bounty_id, COALESCE(json_agg(bhrc), '[]'::json) hunters
		FROM bounty_hunter_rels_cte bhrc
		GROUP BY bhrc.bounty_id
	),res AS (
	    SELECT bc.*, COALESCE(bhrac.hunters, '[]'::json) hunters
	    FROM bounties_cte bc
	    LEFT JOIN bounty_hunter_rels_aggregate_cte bhrac ON bc.id = bhrac.bounty_id
	)
	SELECT
		COALESCE(json_agg(r.*), '[]'::json)
	FROM res r;`
	query, args := util.PgMapQueryV2(query, plan.Params)

	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &util.PgJsonScanWrap{m}, query, args...)
	})
	return
}

func (c *bounties) GetBounty(ctx context.Context, id flake.ID, isOwner bool, output interface{}) (err error) {
	plan := &dbquery.Plan{}
	plan.AddCond(`AND b.id = ${id}`)

	plan.Params = map[string]interface{}{
		"{id}": id,
	}
	var v struct{ Id flake.ID }
	_, err = c.Query(ctx, 0, isOwner, &util.PgJsonScanWrapValues{&[]interface{}{&v}, &[]interface{}{output}}, plan)
	if err != nil {
		return
	}
	if v.Id == 0 {
		err = apierror.ErrNotFound
		return
	}
	return
}

func (c *bounties) CreateUndertakeBounty(ctx context.Context, uid flake.ID, input *coresSdk.CreateUndertakeBountyInput, output *coresSdk.UndertakeBountyResult) (err error) {
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
