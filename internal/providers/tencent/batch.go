package tencent

import (
	"context"
	"sync"

	"github.com/ceheng-io/stock-go/internal/core"
	"github.com/ceheng-io/stock-go/types"
	"github.com/ceheng-io/stock-go/utils"
)

const (
	defaultBatchSize   = 500
	maxBatchSize       = 500
	defaultConcurrency = 7
)

// BatchOptions configures batch quote fetching.
type BatchOptions struct {
	BatchSize   int
	Concurrency int
	OnProgress  func(completed, total int)
}

func normalizeBatchOptions(options BatchOptions) (BatchOptions, error) {
	if options.BatchSize == 0 {
		options.BatchSize = defaultBatchSize
	}
	if options.Concurrency == 0 {
		options.Concurrency = defaultConcurrency
	}
	if options.BatchSize < 0 {
		return BatchOptions{}, invalidArgumentError("batchSize must be positive")
	}
	if options.Concurrency < 0 {
		return BatchOptions{}, invalidArgumentError("concurrency must be positive")
	}
	if options.BatchSize > maxBatchSize {
		options.BatchSize = maxBatchSize
	}
	return options, nil
}

func invalidArgumentError(message string) error {
	return core.NewCodedError("INVALID_ARGUMENT", message, nil)
}

// GetAllQuotesByCodes fetches detailed CN quotes in batches.
func GetAllQuotesByCodes(ctx context.Context, client QuoteClient, codes []string, options BatchOptions) ([]types.FullQuote, error) {
	return batchFetch(ctx, codes, options, func(ctx context.Context, chunk []string) ([]types.FullQuote, error) {
		return GetFullQuotes(ctx, client, chunk)
	})
}

// GetAllHKQuotesByCodes fetches HK quotes in batches.
func GetAllHKQuotesByCodes(ctx context.Context, client QuoteClient, codes []string, options BatchOptions) ([]types.HKQuote, error) {
	return batchFetch(ctx, codes, options, func(ctx context.Context, chunk []string) ([]types.HKQuote, error) {
		return GetHKQuotes(ctx, client, chunk)
	})
}

// GetAllUSQuotesByCodes fetches US quotes in batches.
func GetAllUSQuotesByCodes(ctx context.Context, client QuoteClient, codes []string, options BatchOptions) ([]types.USQuote, error) {
	return batchFetch(ctx, codes, options, func(ctx context.Context, chunk []string) ([]types.USQuote, error) {
		return GetUSQuotes(ctx, client, chunk)
	})
}

func batchFetch[T any](
	ctx context.Context,
	codes []string,
	options BatchOptions,
	fetch func(context.Context, []string) ([]T, error),
) ([]T, error) {
	options, err := normalizeBatchOptions(options)
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		return []T{}, nil
	}

	chunks, err := utils.ChunkArray(codes, options.BatchSize)
	if err != nil {
		return nil, err
	}
	results := make([][]T, len(chunks))
	jobs := make(chan int)
	var completed int
	var completedMu sync.Mutex
	var firstErr error
	var errMu sync.Mutex
	var wg sync.WaitGroup

	workerCount := options.Concurrency
	if workerCount > len(chunks) {
		workerCount = len(chunks)
	}
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				errMu.Lock()
				hasErr := firstErr != nil
				errMu.Unlock()
				if hasErr {
					continue
				}
				value, err := fetch(ctx, chunks[index])
				if err != nil {
					errMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					errMu.Unlock()
					continue
				}
				results[index] = value
				completedMu.Lock()
				completed++
				if options.OnProgress != nil {
					options.OnProgress(completed, len(chunks))
				}
				completedMu.Unlock()
			}
		}()
	}

	for i := range chunks {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	flat := []T{}
	for _, part := range results {
		flat = append(flat, part...)
	}
	return flat, nil
}
