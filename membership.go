package flux

import (
	"github.com/hashicorp/memberlist"
)

func init() {
	config := memberlist.DefaultLocalConfig()
	config.Events = flux

	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}
	list
}

type event struct {
	Flux
}

func (flux *Flux) NotifyJoin(node *memberlist.Node) {
	flux.peers.Add(node.Address())
}

func (flux *Flux) NotifyLeave(node *memberlist.Node) {
	flux.peers.Remove(node.Address)
}

func (flux *Flux) NotifyUpdate(node *memberlist.Node) {

}

func Join(address string) error {
	_, err := list.Join([]string{address})

	return err
}
