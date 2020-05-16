package categories

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/dbquery"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
)

var Categories = &categories{
	Connector: models.DefaultConnector,
}

// categories represents controller for 'categories'.
type categories struct {
	dbconn.Connector
}

// List Categories interface{} by input
func (c *categories) List(ctx context.Context, input *coresSdk.ListCategoriesInput, outputs interface{}) (total int, err error) {
	plan := &dbquery.Plan{}
	plan.RetTotal = true

	if !input.ShowDeleted {
		plan.AddCond(`AND t.deleted = false`)
	}

	if input.Source != "" {
		plan.AddCond(`AND t.source = ${source}`)
	}

	plan.OrderBySql = ` ORDER BY t.created_at DESC`
	plan.LimitSql = ` LIMIT ${limit} OFFSET ${offset}`

	plan.Params = map[string]interface{}{
		"{source}": input.Source,
		"{offset}": input.Offset,
		"{limit}":  input.GetLimit(),
	}

	total, err = c.Query(ctx, outputs, plan)
	return
}

func (c *categories) Get(ctx context.Context, id flake.ID, output interface{}) (err error) {
	plan := &dbquery.Plan{}
	plan.AddCond(`AND t.id = ${id}`)

	plan.Params = map[string]interface{}{
		"{id}": id,
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

func (c *categories) Query(ctx context.Context, m interface{}, plan *dbquery.Plan) (total int, err error) {
	filterSql := `
	FROM categories t
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
		SELECT t.* 
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

func (c *categories) Exists(ctx context.Context, id flake.ID) (exists bool, err error) {
	stmt := "SELECT EXISTS(SELECT 1 FROM categories WHERE  id = ${id})"

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}": id,
	})
	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &exists, query, args...)
	})
	return
}
