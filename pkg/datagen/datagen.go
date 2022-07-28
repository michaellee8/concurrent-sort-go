package datagen

import (
	"context"
	"sync"
)

type DataGenerator struct {
}

func (_ DataGenerator) GenerateData(
	ctx context.Context,
	targetDir string,
	numFiles int,
) (err error) {
	const numOfInt64PerFile = 128 * 1 << 20
	wg := sync.WaitGroup{}
	for fileIdx := 0; fileIdx < numFiles; fileIdx++ {
		currentFileIdx := fileIdx
		go func() {
			nums := make([]int64, 0, numOfInt64PerFile)
			for i := 0; i < numOfInt64PerFile; i++ {
				nums = append(nums, int64(currentFileIdx+i*numFiles))
			}
		}()
	}
	wg.Done()
}
