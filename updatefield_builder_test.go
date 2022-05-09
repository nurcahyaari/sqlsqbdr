package sqlsqbdr_test

import (
	"database/sql"
	"testing"

	"github.com/guregu/null"
	"github.com/nurcahyaari/sqlsqbdr"
	"github.com/stretchr/testify/assert"
)

func TestBuildUpdatedField(t *testing.T) {

	testCase := []struct {
		name string
		exp  func() sqlsqbdr.UpdatedField
		act  func() (sqlsqbdr.UpdatedField, error)
	}{
		{
			name: "test1",
			exp: func() sqlsqbdr.UpdatedField {
				return sqlsqbdr.UpdatedField{
					Name:  []string{"name = ?", "age = ?"},
					Value: []any{"test", 0},
				}
			},
			act: func() (sqlsqbdr.UpdatedField, error) {
				human := struct {
					Name   string `json:"name" db:"name"`
					Age    int    `json:"age" db:"age"`
					Gender string `json:"gender" db:"gender"`
				}{
					Name: "test",
					Age:  0,
				}

				return sqlsqbdr.BuildUpdatedField(human, sqlsqbdr.IncludeField, "name", "age")
			},
		},
		{
			name: "test2 - with null type",
			exp: func() sqlsqbdr.UpdatedField {
				return sqlsqbdr.UpdatedField{
					Name:  []string{"name = ?", "age = ?"},
					Value: []any{"test", null.Int{sql.NullInt64{Valid: true, Int64: 1}}},
				}
			},
			act: func() (sqlsqbdr.UpdatedField, error) {
				human := struct {
					Name   string   `json:"name" db:"name"`
					Age    null.Int `json:"age" db:"age"`
					Gender string   `json:"gender" db:"gender"`
				}{
					Name: "test",
					Age:  null.IntFrom(1),
				}

				return sqlsqbdr.BuildUpdatedField(human, sqlsqbdr.IncludeField, "name", "age")
			},
		},
		{
			name: "test3 - defined struct",
			exp: func() sqlsqbdr.UpdatedField {
				return sqlsqbdr.UpdatedField{
					Name:  []string{"name = ?", "age = ?", "address = ?"},
					Value: []any{"test", 0, "Test1"},
				}
			},
			act: func() (sqlsqbdr.UpdatedField, error) {
				type Human struct {
					Name    string `json:"name" db:"name"`
					Age     int    `json:"age" db:"age"`
					Gender  string `json:"gender" db:"gender"`
					Address string `json:"address" db:"address"`
				}
				human := Human{
					Name:    "test",
					Age:     0,
					Address: "Test1",
				}

				return sqlsqbdr.BuildUpdatedField(human, sqlsqbdr.IncludeField, "name", "age", "address")
			},
		},
		{
			name: "test4 - empty field",
			exp: func() sqlsqbdr.UpdatedField {
				return sqlsqbdr.UpdatedField{
					Name:  []string{"name = ?", "age = ?", "gender = ?", "address = ?"},
					Value: []any{"test", 0, "", "Test1"},
				}
			},
			act: func() (sqlsqbdr.UpdatedField, error) {
				type Human struct {
					Name    string `json:"name" db:"name"`
					Age     int    `json:"age" db:"age"`
					Gender  string `json:"gender" db:"gender"`
					Address string `json:"address" db:"address"`
				}
				human := Human{
					Name:    "test",
					Age:     0,
					Address: "Test1",
				}

				return sqlsqbdr.BuildUpdatedField(human, sqlsqbdr.IncludeField)
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.exp()
			act, err := tc.act()
			assert.NoError(t, err)
			assert.Equal(t, exp, act)
		})
	}
}
