package repository

import (
	"context"
	"errors"

	"github.com/uptrace/bun"
)

var ErrNotFound = errors.New("not found")

type Filters struct {
	Search        string
	SearchColumns []string
	DariTanggal   string
	SampaiTanggal string
	Tahun         int
	DariColumn    string
	SampaiColumn  string
	TahunColumn   string
	Status        string
	StatusColumn  string
	ExactID       int64
	ExactColumn   string
}

func (f Filters) apply(q *bun.SelectQuery) *bun.SelectQuery {
	if f.Search != "" && len(f.SearchColumns) > 0 {
		like := "%" + f.Search + "%"
		args := make([]any, 0, len(f.SearchColumns))
		clause := "("
		for i, col := range f.SearchColumns {
			if i > 0 {
				clause += " OR "
			}
			clause += col + " ILIKE ?"
			args = append(args, like)
		}
		clause += ")"
		q = q.Where(clause, args...)
	}
	if f.DariTanggal != "" && f.DariColumn != "" {
		q = q.Where(f.DariColumn+" >= ?", f.DariTanggal)
	}
	if f.SampaiTanggal != "" && f.SampaiColumn != "" {
		q = q.Where(f.SampaiColumn+" <= ?", f.SampaiTanggal)
	}
	if f.Tahun > 0 && f.TahunColumn != "" {
		q = q.Where("EXTRACT(YEAR FROM "+f.TahunColumn+") = ?", f.Tahun)
	}
	if f.Status != "" && f.StatusColumn != "" {
		q = q.Where(f.StatusColumn+" = ?", f.Status)
	}
	if f.ExactID > 0 && f.ExactColumn != "" {
		q = q.Where(f.ExactColumn+" = ?", f.ExactID)
	}
	return q
}

func paginate[T any](ctx context.Context, base func() *bun.SelectQuery, page, perPage int, order string, dest *[]*T) (int, error) {
	total, err := base().Count(ctx)
	if err != nil {
		return 0, err
	}

	q := base()
	if order != "" {
		q = q.Order(order)
	}
	if err := q.Limit(perPage).Offset((page-1)*perPage).Scan(ctx, dest); err != nil {
		return 0, err
	}
	return total, nil
}
