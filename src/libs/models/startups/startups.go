package startups

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/dbquery"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models"
	coresSdk "cos-backend-com/src/libs/sdk/cores"

	"github.com/jmoiron/sqlx"
)

var StartUps = &startUps{
	Connector: models.DefaultConnector,
}

// startUps represents controller for 'startups'.
type startUps struct {
	dbconn.Connector
}

func (c *startUps) List(ctx context.Context, uid flake.ID, input *coresSdk.ListStartUpsInput, outputs interface{}) (total int, err error) {
	plan := &dbquery.Plan{}
	plan.RetTotal = true

	if uid != 0 {
		plan.AddCond(`AND t.uid = ${uid}`)
	}
	if input.CategoryId != flake.ID(0) {
		plan.AddCond(`AND t.category_id = ${categoryId}`)
	}
	if input.IsIRO {
		plan.AddCond(`AND t.is_iro = ${isIRO}`)
	}
	var keyword string
	if input.Keyword != "" {
		keyword = "%" + util.PgEscapeLike(input.Keyword) + "%"
		plan.AddCond(`AND t.name ILIKE ${keyword}`)
	}

	plan.OrderBySql = ` ORDER BY t.created_at DESC`
	plan.LimitSql = ` LIMIT ${limit} OFFSET ${offset}`

	plan.Params = map[string]interface{}{
		"{uid}":        uid,
		"{categoryId}": input.CategoryId,
		"{isIRO}":      input.IsIRO,
		"{keyword}":    keyword,
		"{offset}":     input.Offset,
		"{limit}":      input.GetLimit(),
	}

	total, err = c.Query(ctx, outputs, plan)
	return
}

func (c *startUps) Get(ctx context.Context, uid, id flake.ID, output interface{}) (err error) {
	plan := &dbquery.Plan{}
	plan.AddCond(`AND t.id = ${id}`)
	plan.AddCond(`AND t.uid = ${uid}`)

	plan.Params = map[string]interface{}{
		"{id}":  id,
		"{uid}": uid,
	}
	var v struct{ Id flake.ID }
	_, err = c.Query(ctx, &util.PgJsonScanWrapValues{&[]interface{}{&v}, &[]interface{}{output}}, plan)
	if err != nil {
		return
	}
	if v.Id == 0 {
		err = apierror.ErrNotFound
		return
	}
	return
}

func (c *startUps) Create(ctx context.Context, uid flake.ID, input *coresSdk.CreateStartUpsInput, output interface{}) error {
	stmt := `
	INSERT INTO startups (
		uid, name, mission, logo, tx_id, description_addr, category_id
	) VALUES (
		${uid}, ${name}, ${mission}, ${logo}, ${txId}, ${descriptionAddr}, ${categoryId}
	) RETURNING id`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{uid}":             uid,
		"{name}":            input.Name,
		"{mission}":         input.Mission,
		"{logo}":            input.Logo,
		"{txId}":            input.TxId,
		"{descriptionAddr}": input.DescriptionAddr,
		"{categoryId}":      input.CategoryId,
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		var id flake.ID
		if er = db.GetContext(ctx, &id, query, args...); er != nil {
			return
		}

		newCtx := dbconn.WithDB(ctx, db)
		er = c.Get(newCtx, uid, id, output)
		return
	})
}

//func (c *startUps) Update(ctx context.Context, enterpriseId, id flake.ID, input *coresSdk.StartUpsInput, output interface{}) error {
//	stmt := `
//	UPDATE startups SET (
//		name, mission, logo, tx_id, blocknum, description_addr, category_id, state, isiro, updated_at
//	) = (
//		${name}, ${mission}, ${logo}, ${txId}, ${blocknum}, ${descriptionAddr}, ${categoryId}, ${state}, ${isiro}, CURRENT_TIMESTAMP
//	) WHERE
//		enterprise_id = ${enterpriseId} AND id = ${id}
//	RETURNING id`
//
//	query, args := util.PgMapQuery(stmt, map[string]interface{}{
//		"{id}":              id,
//		"{enterpriseId}":    enterpriseId,
//		"{name}":            input.Name,
//		"{mission}":         input.Mission,
//		"{logo}":            input.Logo,
//		"{txId}":            input.TxID,
//		"{blocknum}":        input.Blocknum,
//		"{descriptionAddr}": input.DescriptionAddr,
//		"{categoryId}":      input.CategoryID,
//		"{state}":           input.State,
//		"{isiro}":           input.Isiro,
//	})
//
//	return c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
//		var id flake.ID
//		if er = db.GetContext(ctx, &id, query, args...); er != nil {
//			return
//		}
//
//		newCtx := dbconn.WithDB(ctx, db)
//		er = c.Get(newCtx, enterpriseId, id, output)
//		return
//	})
//}

func (c *startUps) Query(ctx context.Context, m interface{}, plan *dbquery.Plan) (total int, err error) {
	filterSql := `
	FROM startups t
		INNER JOIN categories c ON c.id = t.category_id
	WHERE 1=1` + plan.Conditions

	if plan.RetTotal {
		query := `SELECT count(DISTINCT t.id) ` + filterSql
		query, args := util.PgMapQueryV2(query, plan.Params)

		err = c.Invoke(ctx, func(db dbconn.Q) error {
			return db.GetContext(ctx, &total, query, args...)
		})
		if err != nil {
			return
		}
	}

	query := `
	WITH res AS (
		SELECT t.*,
               c as category
		` + filterSql + `
		` + plan.OrderBySql + `
		` + plan.LimitSql + `
	)
	SELECT
		COALESCE(json_agg(r.*), '[]'::json)
	FROM res r`
	query, args := util.PgMapQueryV2(query, plan.Params)

	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &util.PgJsonScanWrap{m}, query, args...)
	})
	return
}

func (c *startUps) Exists(ctx context.Context, uid, id flake.ID) (exists bool, err error) {
	stmt := "SELECT EXISTS(SELECT 1 FROM startups WHERE uid = ${uid} AND id = ${id})"

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":           id,
		"{enterpriseId}": uid,
	})
	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &exists, query, args...)
	})
	return
}
