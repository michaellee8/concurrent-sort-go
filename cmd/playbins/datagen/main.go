package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, func(), error) {
	prodConfig := zap.NewProductionConfig()
	prodConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, err := prodConfig.Build()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		_ = logger.Sync()
	}
	return logger, cleanup, nil
}

func main() {
	var (
		targetDir string
		numFiles  int
	)
	flag.IntVar(&numFiles, "numFiles", 0, "number of files")
	flag.StringVar(&targetDir, "targetDir", "", "target directory")
	flag.Parse()
	if numFiles == 0 || targetDir == "" {
		flag.PrintDefaults()
		return
	}
	dg, cleanup, err := InitializeDataGenerator()
	if err != nil {
		panic(err)
	}
	defer cleanup()
	err = dg.GenerateData(context.TODO(), targetDir, numFiles)
	if err != nil {
		panic(err)
	}
}
