//go:build wireinject
// +build wireinject

package users

import (
	"github.com/google/wire"
	"github.com/xans-me/authopia/app"
)

var (
	ModuleSet = wire.NewSet(
		NewRepository,
		NewUseCase,
		wire.Bind(new(IUseCase), new(*UseCase)),
		wire.Bind(new(IUsersRepository), new(*Repository)),
		NewRpcDelivery)
)

func InjectRPC() (*RpcDelivery, error) {
	panic(wire.Build(app.AppModule, ModuleSet))
}
