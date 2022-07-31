package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
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
	dg, err := InitializeDataGenerator()
	if err != nil {
		panic(err)
	}
	err = dg.GenerateData(context.TODO(), targetDir, numFiles)
	if err != nil {
		panic(err)
	}
}
