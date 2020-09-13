package huntersModels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
)

var Hunters = &hunters{
	Connector: models.DefaultConnector,
}

type hunters struct {
	dbconn.Connector
}

func (c *hunters) Get(ctx context.Context, id flake.ID, output interface{}) (err error) {
	stmt := `
		WITH res AS (
			SELECT *
			FROM hunters
			WHERE id = ${id}
		)SELECT row_to_json(r) FROM res r;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}": id,
	})
	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{&output}, query, args...)
	})
}

func (c *hunters) Update(ctx context.Context, uid flake.ID, input *coresSdk.UpdateHunterInput) (err error) {
	stmt := `
		INSERT INTO hunters(user_id, name, skills, about, description_addr, email)
		VALUES (
		    ${userId}, ${name}, ARRAY [${skills}], ${about}, ${descriptionAddr}, ${email}
		)
		ON CONFLICT DO
		UPDATE SET (
		    name, skills, about, description_addr, email
		) = (
		    ${name}, ARRAY [${skills}], ${about},${descriptionAddr}, ${email}
		);
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{userId}":          uid,
		"{name}":            input.Name,
		"{skills}":          input.Skills,
		"{about}":           input.About,
		"{descriptionAddr}": input.DescriptionAddr,
		"{email}":           input.Email,
	})
	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		_, er = db.ExecContext(ctx, query, args...)
		return er
	})
}
