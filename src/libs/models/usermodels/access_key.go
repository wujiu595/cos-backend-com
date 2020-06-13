package usermodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
)

var AccessKeys = &accessKeys{models.DefaultConnector}

type accessKeys struct {
	dbconn.Connector
}

func (p *accessKeys) FindByAccessKey(ctx context.Context, accessKey string, m interface{}) (err error) {
	query := `
	SELECT row_to_json(a) FROM access_keys a WHERE a.key = ${key}
`

	query, args := util.PgMapQueryV2(query, map[string]interface{}{
		"{key}": accessKey,
	})
	err = p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &util.PgJsonScanWrap{m}, query, args...)
	})
	return
}
