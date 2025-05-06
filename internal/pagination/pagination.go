package pagination

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
)

type Paginator struct {
	defaultLimit, defaultOffset uint
}

func NewPaginator(defaultLimit, defaultOffset uint) *Paginator {
	if defaultLimit == 0 {
		defaultLimit = 1
	}
	return &Paginator{defaultLimit: defaultLimit, defaultOffset: defaultOffset}
}

func (p *Paginator) LimitFromQueryParams(queryParams url.Values) uint {
	const op = "[internal/pagination LimitFromQueryParams]"

	limit := queryParams.Get("limit")
	limitUint, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		log.Warnf("%s invalid `limit` parameter: `%s`", op, limit)
		return p.defaultLimit
	}
	return uint(limitUint)
}

func (p *Paginator) OffsetFromQueryParams(queryParams url.Values) uint {
	const op = "[internal/pagination OffsetFromQueryParams]"

	offset := queryParams.Get("offset")
	offsetUint, err := strconv.ParseUint(offset, 10, 64)
	if err != nil {
		log.Warnf("%s invalid `offset` parameter: `%s`", op, offset)
		return p.defaultOffset
	}
	return uint(offsetUint)
}
