package store

import (
	sq "github.com/Masterminds/squirrel"
)

const maxPageSize = 500

func withPagination(pageSize, pageNumber int) func(sq.SelectBuilder) sq.SelectBuilder {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		if pageSize > maxPageSize {
			pageSize = maxPageSize
		}
		offset := (pageNumber - 1) * pageSize
		return builder.Limit(uint64(pageSize)).Offset(uint64(offset))
	}
}
