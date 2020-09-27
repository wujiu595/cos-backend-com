package huntersModels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	accountSdk "cos-backend-com/src/libs/sdk/account"
)

var Hunters = &hunters{
	Connector: models.DefaultConnector,
}

type hunters struct {
	dbconn.Connector
}

func (c *hunters) Upsert(ctx context.Context, uid flake.ID, input *accountSdk.UpdateHunterInput) (err error) {
	stmt := `
		INSERT INTO hunters(user_id, name, skills, about, description_addr, email)
		VALUES (
		    ${userId}, ${name}, ARRAY [${skills}], ${about}, ${descriptionAddr}, ${email}
		)
		ON CONFLICT(user_id) DO
		UPDATE SET (
		    name, skills, about, description_addr, email
		) = (
		    ${name}, ARRAY [${skills}], ${about},${descriptionAddr}, ${email}
		);
	`
	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
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
