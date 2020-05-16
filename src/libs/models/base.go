package models

import (
	"cos-backend-com/src/common/dbconn"
)

var DefaultConnector dbconn.Connector = dbconn.RegisterConnector(dbconn.Connectors.Default)
