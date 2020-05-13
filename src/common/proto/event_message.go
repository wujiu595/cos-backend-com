package proto

import (
	"time"

	"cos-backend-com/src/common/flake"
)

const PreAggTopic = "pre_agg_event"

type EventMessage struct {
	ComponentId flake.ID  `json:"componentId"`
	ProcessId   flake.ID  `json:"processId"`
	DimensionId flake.ID  `json:"dimensionId"`
	Time        time.Time `json:"time"`
}
