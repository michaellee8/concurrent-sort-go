//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/michaellee8/concurrent-sort-go/pkg/datagen"
)

func InitializeDataGenerator() (*datagen.DataGenerator, func(), error) {
	wire.Build(
		NewLogger,
		datagen.NewDataGenerator,
	)
	return &datagen.DataGenerator{}, nil, nil
}
