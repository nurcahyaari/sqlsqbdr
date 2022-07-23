package sqlsqbdr_test

import (
	"testing"

	"github.com/nurcahyaari/sqlsqbdr"
	"github.com/stretchr/testify/assert"
)

func TestBuildWhereFilter(t *testing.T) {
	testCase := []struct {
		name string
		exp  func() (string, []interface{})
		act  func() (string, []interface{})
	}{
		{
			name: "test1",
			exp: func() (string, []interface{}) {
				return "name = ?", []interface{}{"test"}
			},
			act: func() (string, []interface{}) {
				return sqlsqbdr.BuildWhereFilter(sqlsqbdr.Filters{
					&sqlsqbdr.Filter{
						Field: "name",
						Value: "test",
					},
				})
			},
		},
		{
			name: "test2 - Multi filter",
			exp: func() (string, []interface{}) {
				return "name = ? AND age = ?", []interface{}{"test", 1}
			},
			act: func() (string, []interface{}) {
				return sqlsqbdr.BuildWhereFilter(sqlsqbdr.Filters{
					&sqlsqbdr.Filter{
						Field: "name",
						Value: "test",
					},
					&sqlsqbdr.Filter{
						Field: "age",
						Value: 1,
					},
				})
			},
		},
		{
			name: "test2 - IN AND NOT IN",
			exp: func() (string, []interface{}) {
				return "name IN ? OR age NOT IN ?", []interface{}{[]string{"test", "test1"}, []int64{1, 2}}
			},
			act: func() (string, []interface{}) {
				return sqlsqbdr.BuildWhereFilter(sqlsqbdr.Filters{
					&sqlsqbdr.Filter{
						Field:                 "name",
						Value:                 []string{"test", "test1"},
						ComparisonOperator:    sqlsqbdr.FilterIn,
						ConcatenationOperator: sqlsqbdr.FilterOr,
					},
					&sqlsqbdr.Filter{
						Field:              "age",
						Value:              []int64{1, 2},
						ComparisonOperator: sqlsqbdr.FilterNotIn,
					},
				})
			},
		},
		{
			name: "test3 - raw (between)",
			exp: func() (string, []interface{}) {
				return "age between ? AND ?", []interface{}{[]int{1, 2}}
			},
			act: func() (string, []interface{}) {
				return sqlsqbdr.BuildWhereFilter(sqlsqbdr.Filters{
					&sqlsqbdr.Filter{
						Field:              "age between ? AND ?",
						Value:              []int{1, 2},
						ComparisonOperator: sqlsqbdr.FilterRaw,
					},
				})
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			expQ, expV := tc.exp()
			actQ, actV := tc.act()
			assert.Equal(t, expQ, actQ)
			assert.Equal(t, expV, actV)
		})
	}
}
