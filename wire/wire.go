//go:build wireinject

package main

import (
	"github.com/google/wire"
)

var providerSet wire.ProviderSet = wire.NewSet(
	//NewMessage,
	wire.Struct(new(Message), "Content", "Code"), // fieldNames 确定哪些字段可以被赋值
	//wire.Struct(new(Message), "*"),
	NewGreeter,
)

func InitializeEvent(phrase string, code int) (Event, error) {
	wire.Build(NewEvent, providerSet)
	return Event{}, nil
}
