package stock

import (
	"context"

	"github.com/ceheng.io/stock-go/utils"
)

// ChunkArray 将切片分割成指定大小的块。
func ChunkArray[T any](values []T, chunkSize int) ([][]T, error) {
	return utils.ChunkArray(values, chunkSize)
}

// AsyncPool 按指定并发数执行任务；preserveOrder 为 true 时按任务顺序返回结果。
func AsyncPool[T any](
	ctx context.Context,
	tasks []func(context.Context) (T, error),
	concurrency int,
	preserveOrder bool,
) ([]T, error) {
	return utils.AsyncPool(ctx, tasks, concurrency, preserveOrder)
}

// AssertPositiveInteger 校验值为正整数。
func AssertPositiveInteger(value int, name string) error {
	return utils.AssertPositiveInteger(value, name)
}

// AssertKlinePeriod 校验历史 K 线周期。
func AssertKlinePeriod(period string) error {
	return utils.AssertKlinePeriod(period)
}

// AssertMinutePeriod 校验分钟 K 线周期。
func AssertMinutePeriod(period string) error {
	return utils.AssertMinutePeriod(period)
}

// AssertAdjustType 校验复权类型。
func AssertAdjustType(adjust string) error {
	return utils.AssertAdjustType(adjust)
}

// PeriodCode 返回东方财富历史 K 线周期代码。
func PeriodCode(period string) (string, error) {
	return utils.PeriodCode(period)
}

// GetPeriodCode 返回东方财富历史 K 线周期代码。
func GetPeriodCode(period string) (string, error) {
	return PeriodCode(period)
}

// AdjustCode 返回东方财富复权类型代码。
func AdjustCode(adjust string) (string, error) {
	return utils.AdjustCode(adjust)
}

// GetAdjustCode 返回东方财富复权类型代码。
func GetAdjustCode(adjust string) (string, error) {
	return AdjustCode(adjust)
}
