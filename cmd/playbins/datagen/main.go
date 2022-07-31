package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

func NewLogger() (*zap.Logger, func(), error) {
	logDirPath := os.Getenv("LOG_DIR")
	prodConfig := zap.NewProductionConfig()
	prodConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	if logDirPath != "" {
		prodConfig.OutputPaths = append(
			prodConfig.OutputPaths,
			filepath.Join(logDirPath, fmt.Sprintf("concurrent-sort-go-%d.log.txt", time.Now().Unix())),
		)
	}
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
	defer profile.Start(profile.MemProfileHeap).Stop()
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
