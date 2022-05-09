package sqlsqbdr_test

import (
	"fmt"
	"testing"

	"github.com/nurcahyaari/sqlsqbdr"
	"github.com/stretchr/testify/assert"
)

func TestBuildWhereFilter(t *testing.T) {
	testCase := []struct {
		name string
		exp  func() string
		act  func() string
	}{
		{
			name: "test1",
			exp: func() string {
				return "name = \"test\""
			},
			act: func() string {
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
			exp: func() string {
				return "name = \"test\" AND age = 1"
			},
			act: func() string {
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
			exp: func() string {
				return "name IN (\"test\",\"test1\") OR age NOT IN (1,2)"
			},
			act: func() string {
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
			exp: func() string {
				return "age between 1 AND 2"
			},
			act: func() string {
				return sqlsqbdr.BuildWhereFilter(sqlsqbdr.Filters{
					&sqlsqbdr.Filter{
						Value:              fmt.Sprintf("%s %s %v AND %v", "age", "between", 1, 2),
						ComparisonOperator: sqlsqbdr.FilterRaw,
					},
				})
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.exp()
			act := tc.act()
			assert.Equal(t, exp, act)
		})
	}
}
