// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/michaellee8/concurrent-sort-go/pkg/datagen"
)

// Injectors from wire.go:

func InitializeDataGenerator() (*datagen.DataGenerator, error) {
	logger, err := NewLogger()
	if err != nil {
		return nil, err
	}
	dataGenerator := datagen.NewDataGenerator(logger)
	return dataGenerator, nil
}
