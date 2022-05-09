package sqlsqbdr

import (
	"fmt"
	"reflect"
	"strings"
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

func BuildWhereFilter(filters Filters) string {
	filterQuery := ""

	for i, r := range filters {
		typeOf := reflect.TypeOf(r.Value)
		valueOf := reflect.ValueOf(r.Value)

		var listvalue []string
		if typeOf.Kind() == reflect.Array || typeOf.Kind() == reflect.Slice {
			for i := 0; i < valueOf.Len(); i++ {
				valueOfIndex := valueOf.Index(i)
				if valueOfIndex.Kind() == reflect.String {
					valueOfIndex.SetString(fmt.Sprintf("\"%s\"", valueOfIndex.String()))
				} else {
					valueOfIndex = reflect.ValueOf(fmt.Sprintf("%v", valueOfIndex.Interface()))
				}
				listvalue = append(listvalue, valueOfIndex.String())
			}
		}

		if r.ComparisonOperator.IsMultiValue() {
			r.Value = fmt.Sprintf("(%v)", strings.Join(listvalue, ","))
		}

		if typeOf.Kind() == reflect.String {
			r.Value = fmt.Sprintf("\"%s\"", r.Value)
		}

		if r.ComparisonOperator.IsRaw() {
			filterQuery += valueOf.String()
		} else {
			filterQuery += fmt.Sprintf("%s %s %v", r.Field, r.ComparisonOperator.ToString(), r.Value)
		}

		if i < len(filters)-1 {
			filterQuery += fmt.Sprintf(" %s ", r.ConcatenationOperator.ToString())
		}
	}

	return filterQuery
}
