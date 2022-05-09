package sqlsqbdr_test

import (
	"fmt"
	"testing"

	"github.com/nurcahyaari/sqlsqbdr"
	"github.com/stretchr/testify/assert"
)

func TestInsertFieldBuilder(t *testing.T) {
	testCase := []struct {
		name     string
		expected func() sqlsqbdr.InsertField
		actual   func() (sqlsqbdr.InsertField, error)
	}{
		{
			name: "test1",
			expected: func() sqlsqbdr.InsertField {
				return sqlsqbdr.InsertField{
					Name:        []string{"name", "age", "gender"},
					Placeholder: []string{"(?,?,?)", "(?,?,?)"},
					Values:      [][]any{{"test", 1, "M"}, {"test2", 2, "M"}},
				}
			},
			actual: func() (sqlsqbdr.InsertField, error) {
				humans := []struct {
					Name   string `json:"name" db:"name"`
					Age    int    `json:"age" db:"age"`
					Gender string `json:"gender" db:"gender"`
				}{
					{
						Name:   "test",
						Age:    1,
						Gender: "M",
					},
					{
						Name:   "test2",
						Age:    2,
						Gender: "M",
					},
				}
				return sqlsqbdr.BuildInsertField(humans, sqlsqbdr.IncludeField)
			},
		},
		{
			name: "test2 - single",
			expected: func() sqlsqbdr.InsertField {
				return sqlsqbdr.InsertField{
					Name:        []string{"name", "age", "gender"},
					Placeholder: []string{"(?,?,?)"},
					Values:      [][]any{{"test", 1, "M"}},
				}
			},
			actual: func() (sqlsqbdr.InsertField, error) {
				humans := []struct {
					Name   string `json:"name" db:"name"`
					Age    int    `json:"age" db:"age"`
					Gender string `json:"gender" db:"gender"`
				}{
					{
						Name:   "test",
						Age:    1,
						Gender: "M",
					},
				}
				return sqlsqbdr.BuildInsertField(humans, sqlsqbdr.IncludeField)
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			act, err := tc.actual()
			fmt.Println(act)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected(), act)
		})
	}
}
