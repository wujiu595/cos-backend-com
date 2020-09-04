package followmodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
)

var Follows = &follows{
	Connector: models.DefaultConnector,
}

type follows struct {
	dbconn.Connector
}

func (c *follows) CreateFollow(ctx context.Context, startupId, uid flake.ID, id flake.ID) (err error) {
	stmt := `
		INSERT INTO startups_follows_rel(id, startup_id, user_id)
		VALUES (${id}, ${startupId},${uid});
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{uid}":       uid,
		"{id}":        id,
		"{startupId}": startupId,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		_, er = db.ExecContext(ctx, query, args...)
		return er
	})
}

func (c *follows) NextId(ctx context.Context) (netxtId flake.ID, err error) {
	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &netxtId, `SELECT id_generator()`)
	})
	return
}
