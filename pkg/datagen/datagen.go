package datagen

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/michaellee8/concurrent-sort-go/pkg/concsort"
	"github.com/pbnjay/memory"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type DataGenerator struct {
	logger *zap.Logger
}

func NewDataGenerator(logger *zap.Logger) *DataGenerator {
	return &DataGenerator{
		logger: logger,
	}
}

func (dg *DataGenerator) GenerateData(
	ctx context.Context,
	targetDir string,
	numFiles int,
) (err error) {
	const errMsg = "cannot generate data"
	const numOfInt64PerFile = 128 * 1 << 20
	maxConcurrentFiles := int(memory.TotalMemory() / 2 / (8 * numOfInt64PerFile))
	if maxConcurrentFiles > runtime.NumCPU()-1 {
		maxConcurrentFiles = runtime.NumCPU() - 1
	}
	errGp := errgroup.Group{}
	errGp.SetLimit(maxConcurrentFiles)
	dg.logger.Info(
		"GenerateData: Max number of files being generated concurrently",
		zap.Int("maxConcurrentFiles", maxConcurrentFiles),
	)

	// Pre-allocate all memory segment that will be used for shuffling random numbers
	// so that we can avoid thrashing and memory allocation/reallocation.
	numsMemPool := sync.Pool{
		New: func() any {
			dg.logger.Debug("New nums allocated")
			return make([]int64, numOfInt64PerFile)
		},
	}
	for i := 0; i < maxConcurrentFiles; i++ {
		dg.logger.Debug(
			"nums being prepared and put into numsMemPool",
			zap.Int("i", i),
		)
		numsMemPool.Put(make([]int64, numOfInt64PerFile))
	}
	rand.Seed(time.Now().UnixNano())
	for fileIdx := 0; fileIdx < numFiles; fileIdx++ {
		currentFileIdx := fileIdx
		errGp.Go(func() (err error) {
			dg.logger.Debug(
				"file generation started",
				zap.Int("currentFileIdx", currentFileIdx),
			)
			nums := numsMemPool.Get().([]int64)
			// Make sure it is put back after usage.
			defer numsMemPool.Put(nums)
			for i := 0; i < numOfInt64PerFile; i++ {
				nums[i] = int64(currentFileIdx + i*numFiles)
			}
			rand.Shuffle(len(nums), func(i, j int) {
				nums[i], nums[j] = nums[j], nums[i]
			})
			f, err := os.Create(filepath.Join(targetDir, fmt.Sprintf("%d.bin", currentFileIdx)))
			if err != nil {
				return errors.Wrap(err, errMsg)
			}
			defer func() {
				_ = f.Close()
			}()
			for i := range nums {
				if err := binary.Write(f, concsort.DefaultEndian, nums[i]); err != nil {
					return errors.Wrap(err, "cannot write to file")
				}
			}
			dg.logger.Debug(
				"file generation ended",
				zap.Int("currentFileIdx", currentFileIdx),
			)
			return nil
		})
	}
	if err := errGp.Wait(); err != nil {
		return errors.Wrap(err, errMsg)
	}
	return nil
}
