package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/cstdev/moonapi/query"
	log "github.com/sirupsen/logrus"
)

// RequestQuery represents all the fields available in a query but as
// strings.
type RequestQuery struct {
	Order         string
	Asc           string
	Configuration string
	HoldSet       string
	Filter        string
	MinGrade      string
	MaxGrade      string
	Page          string
	PageSize      string
}

// Query takes a RequestQuery and converts it properties to the correct types
// to make a Query object which can then be used to perform a search.
func (q *RequestQuery) Query() (query.Query, error) {
	var asc bool
	var page int
	var pageSize int
	var err error

	if q.Page != "" {
		page, err = strconv.Atoi(q.Page)
		if err != nil {
			return nil, errors.New("Invalid page number")
		}
	}

	if q.PageSize != "" {
		pageSize, err = strconv.Atoi(q.PageSize)
		if err != nil {
			return nil, errors.New("Invalid page size.")
		}
	}

	if q.Asc != "" {
		asc, err = strconv.ParseBool(q.Asc)
		if err != nil {
			return nil, errors.New("Invalid ascending value, should be 'true' or 'false'")
		}
	}

	builder := query.New()
	if q.Order != "" {
		orderType, err := query.ToOrder(q.Order)
		if err != nil {
			return nil, err
		}
		builder.Sort(*orderType, asc)
	}

	if q.Configuration != "" {
		configType, err := query.ToConfiguration(q.Configuration)
		if err != nil {
			return nil, err
		}
		builder.Configuration(*configType)
	}

	if q.HoldSet != "" {
		sets := strings.Split(q.HoldSet, ",")
		for _, set := range sets {
			set = strings.TrimSpace(set)
			log.WithFields(log.Fields{
				"holdSet": set,
			}).Debug("Split hold sets")
			holdType, err := query.ToHoldSet(set)
			if err != nil {
				return nil, err
			}

			builder.HoldSet(*holdType)
		}
	}

	if q.Filter != "" {
		filterType, err := query.ToFilter(q.Filter)
		if err != nil {
			return nil, err
		}
		builder.Filter(*filterType)
	}

	if q.MinGrade != "" {
		minGradeType, err := query.ToGrade(q.MinGrade)
		if err != nil {
			return nil, err
		}

		builder.MinGrade(*minGradeType)
	}

	if q.MaxGrade != "" {
		maxGradeType, err := query.ToGrade(q.MaxGrade)
		if err != nil {
			return nil, err
		}

		builder.MaxGrade(*maxGradeType)
	}

	if page > 0 {
		builder.Page(page)
	}

	if pageSize > 0 {

		builder.PageSize(pageSize)
	}

	log.WithFields(log.Fields{
		"RequestQuery": q,
	}).Debug("Building query.")

	query, errs := builder.Build()
	if err != nil {
		return nil, errs[0]
	}
	return query, nil
}
