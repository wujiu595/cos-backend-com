package tagmodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
)

var Tags = &tags{
	Connector: models.DefaultConnector,
}

// categories represents controller for 'categories'.
type tags struct {
	dbconn.Connector
}

// List Categories interface{} by input
func (c *tags) List(ctx context.Context, input *coresSdk.ListTagsInput, output interface{}) (err error) {
	stmt := `
	WITH res AS (
        SELECT name FROM tags WHERE source = ${source}
	)
	SELECT
		COALESCE(json_build_array(r.name), '[]'::json)
	FROM res r;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{source}": input.Source,
	})

	return c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}
