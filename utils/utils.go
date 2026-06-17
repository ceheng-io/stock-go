package utils

import (
	"context"
	"fmt"
	"sync"

	"github.com/ceheng-io/stock-go/internal/core"
)

var (
	klinePeriods = map[string]string{
		"daily":   "101",
		"weekly":  "102",
		"monthly": "103",
	}
	minutePeriods = map[string]struct{}{
		"1":  {},
		"5":  {},
		"15": {},
		"30": {},
		"60": {},
	}
	adjustTypes = map[string]string{
		"":     "0",
		"none": "0",
		"qfq":  "1",
		"hfq":  "2",
	}
)

// AssertPositiveInteger 校验值为正整数。
func AssertPositiveInteger(value int, name string) error {
	if value <= 0 {
		return invalidArgumentError(fmt.Sprintf("%s must be a positive integer", name))
	}
	return nil
}

// AssertKlinePeriod 校验历史 K 线周期。
func AssertKlinePeriod(period string) error {
	if _, ok := klinePeriods[period]; !ok {
		return invalidArgumentError("period must be one of: daily, weekly, monthly")
	}
	return nil
}

// AssertMinutePeriod 校验分钟 K 线周期。
func AssertMinutePeriod(period string) error {
	if _, ok := minutePeriods[period]; !ok {
		return invalidArgumentError("period must be one of: 1, 5, 15, 30, 60")
	}
	return nil
}

// AssertAdjustType 校验复权类型。
func AssertAdjustType(adjust string) error {
	if _, ok := adjustTypes[adjust]; !ok {
		return invalidArgumentError("adjust must be one of: '', 'none', 'qfq', 'hfq'")
	}
	return nil
}

// ChunkArray 将切片分割成指定大小的块。
func ChunkArray[T any](values []T, chunkSize int) ([][]T, error) {
	if chunkSize <= 0 {
		return nil, invalidArgumentError("chunkSize must be a positive integer")
	}
	if len(values) == 0 {
		return [][]T{}, nil
	}
	chunks := make([][]T, 0, (len(values)+chunkSize-1)/chunkSize)
	for start := 0; start < len(values); start += chunkSize {
		end := start + chunkSize
		if end > len(values) {
			end = len(values)
		}
		chunks = append(chunks, values[start:end])
	}
	return chunks, nil
}

// AsyncPool 按指定并发数执行任务；preserveOrder 为 true 时按任务顺序返回结果。
func AsyncPool[T any](
	ctx context.Context,
	tasks []func(context.Context) (T, error),
	concurrency int,
	preserveOrder bool,
) ([]T, error) {
	if concurrency <= 0 {
		return nil, invalidArgumentError("concurrency must be a positive integer")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return []T{}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	workerCount := concurrency
	if workerCount > len(tasks) {
		workerCount = len(tasks)
	}

	jobs := make(chan int)
	var wg sync.WaitGroup
	var resultMu sync.Mutex
	var errMu sync.Mutex
	var firstErr error

	var ordered []T
	var unordered []T
	if preserveOrder {
		ordered = make([]T, len(tasks))
	} else {
		unordered = make([]T, 0, len(tasks))
	}

	recordErr := func(err error) {
		if err == nil {
			return
		}
		errMu.Lock()
		if firstErr == nil {
			firstErr = err
			cancel()
		}
		errMu.Unlock()
	}

	worker := func() {
		defer wg.Done()
		for index := range jobs {
			if err := ctx.Err(); err != nil {
				recordErr(err)
				continue
			}
			value, err := tasks[index](ctx)
			if err != nil {
				recordErr(err)
				continue
			}
			resultMu.Lock()
			if preserveOrder {
				ordered[index] = value
			} else {
				unordered = append(unordered, value)
			}
			resultMu.Unlock()
		}
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	for i := range tasks {
		if err := ctx.Err(); err != nil {
			recordErr(err)
			break
		}
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	if preserveOrder {
		return ordered, nil
	}
	return unordered, nil
}

// PeriodCode 返回东方财富历史 K 线周期代码。
func PeriodCode(period string) (string, error) {
	code, ok := klinePeriods[period]
	if !ok {
		return "", invalidArgumentError("period must be one of: daily, weekly, monthly")
	}
	return code, nil
}

// AdjustCode 返回东方财富复权类型代码。
func AdjustCode(adjust string) (string, error) {
	code, ok := adjustTypes[adjust]
	if !ok {
		return "", invalidArgumentError("adjust must be one of: '', 'none', 'qfq', 'hfq'")
	}
	return code, nil
}

func invalidArgumentError(message string) error {
	return core.NewCodedError("INVALID_ARGUMENT", message, nil)
}
