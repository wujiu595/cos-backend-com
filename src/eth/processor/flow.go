package processor

import (
	"github.com/trustmaster/goflow"
)

func NewConfirmingApp() *goflow.Graph {
	n := goflow.NewGraph()
	n.Add("confirmer", new(Confirmer))
	n.Add("updater", new(Updater))
	n.Connect("confirmer", "TransactionOutput", "updater", "TransactionInput")
	n.MapInPort("In", "confirmer", "TransactionInput")
	return n
}
