package daemon

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewHandles, NewServer)

func NewHandles() []func() {
	return []func(){}
}
