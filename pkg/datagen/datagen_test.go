package datagen

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func ExampleDataGenerator_GenerateData() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	targetDir := filepath.Join(cwd, "testdata/tmp")
	if err := os.MkdirAll(targetDir, 0750); err != nil {
		panic(err)
	}
	dg := NewDataGenerator(logger)
	err = dg.GenerateData(context.TODO(), targetDir, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println("success")
	// Output: success
}
