package screener

import (
	"fmt"
	"math"
	"sort"

	"github.com/ceheng.io/stock-go/internal/core"
)

// SortDirection controls numeric sorting order.
type SortDirection string

const (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

// Builder is a local chainable screener.
type Builder[T any] struct {
	items []T
	err   error
}

// Screen creates a local screener from items.
func Screen[T any](items []T) *Builder[T] {
	copied := append([]T(nil), items...)
	return &Builder[T]{items: copied}
}

// Where keeps items matching predicate.
func (b *Builder[T]) Where(predicate func(T) bool) *Builder[T] {
	if b.err != nil {
		return b
	}
	filtered := make([]T, 0, len(b.items))
	for _, item := range b.items {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	b.items = filtered
	return b
}

// SortBy sorts by a nullable numeric selector. Invalid values sink to the end.
func (b *Builder[T]) SortBy(selector func(T) *float64, direction SortDirection) *Builder[T] {
	if b.err != nil {
		return b
	}
	sign := 1
	if direction != Asc {
		sign = -1
	}
	sort.SliceStable(b.items, func(i int, j int) bool {
		left, leftOK := finiteValue(selector(b.items[i]))
		right, rightOK := finiteValue(selector(b.items[j]))
		if !leftOK && !rightOK {
			return false
		}
		if !leftOK {
			return false
		}
		if !rightOK {
			return true
		}
		return float64(sign)*(left-right) < 0
	})
	return b
}

// Top returns the first n items.
func (b *Builder[T]) Top(n int) ([]T, error) {
	if b.err != nil {
		return nil, b.err
	}
	if n < 0 {
		return nil, invalidArgumentError(fmt.Sprintf("top(n): n must be a non-negative integer, got %d", n))
	}
	if n > len(b.items) {
		n = len(b.items)
	}
	return append([]T(nil), b.items[:n]...), nil
}

func invalidArgumentError(message string) error {
	return core.NewCodedError("INVALID_ARGUMENT", message, nil)
}

// ToArray returns all current items.
func (b *Builder[T]) ToArray() ([]T, error) {
	if b.err != nil {
		return nil, b.err
	}
	return append([]T(nil), b.items...), nil
}

// Float returns a non-null float pointer.
func Float(value float64) *float64 {
	return &value
}

func finiteValue(value *float64) (float64, bool) {
	if value == nil || math.IsNaN(*value) || math.IsInf(*value, 0) {
		return 0, false
	}
	return *value, true
}

func round(value float64) float64 {
	return math.Round(value*100) / 100
}
