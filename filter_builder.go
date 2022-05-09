package sqlsqbdr

import (
	"fmt"
)

type FilterComparisonOperator int

const (
	FilterEqual FilterComparisonOperator = iota
	FilterLess
	FilterLessOrEqual
	FilterGreater
	FilterGreaterOrEqual
	FilterIn
	FilterNotIn
	FilterNotEqual
	FilterIsNull
	FilterRaw
)

func (f FilterComparisonOperator) ToString() string {
	var s string
	switch f {
	case FilterLess:
		s = "<"
	case FilterLessOrEqual:
		s = "<="
	case FilterGreater:
		s = ">"
	case FilterGreaterOrEqual:
		s = ">="
	case FilterIn:
		s = "IN"
	case FilterNotIn:
		s = "NOT IN"
	case FilterNotEqual:
		s = "!="
	case FilterIsNull:
		s = "IS NULL"
	default:
		s = "="
	}
	return s
}

func (f FilterComparisonOperator) IsRaw() bool {
	return f == FilterRaw
}

func (f FilterComparisonOperator) IsMultiValue() bool {
	return f == FilterIn || f == FilterNotIn
}

type FilterConcatenationOperator int

const (
	FilterAnd FilterConcatenationOperator = iota
	FilterOr
)

func (f FilterConcatenationOperator) ToString() string {
	var s string
	switch f {
	case FilterOr:
		s = "OR"
	default:
		s = "AND"
	}
	return s
}

type Filter struct {
	Field                 string
	Value                 interface{}
	ComparisonOperator    FilterComparisonOperator
	ConcatenationOperator FilterConcatenationOperator
}

type Filters []*Filter

func BuildWhereFilter(filters Filters) (string, []interface{}) {
	filterQuery := ""
	var values []interface{}

	for i, r := range filters {
		values = append(values, r.Value)

		if r.ComparisonOperator.IsRaw() {
			filterQuery += r.Field
		} else {
			filterQuery += fmt.Sprintf("%s %s %v", r.Field, r.ComparisonOperator.ToString(), "?")
		}

		if i < len(filters)-1 {
			filterQuery += fmt.Sprintf(" %s ", r.ConcatenationOperator.ToString())
		}
	}

	return filterQuery, values
}
