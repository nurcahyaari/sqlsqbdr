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
		name    string
		isError bool
		exp     func() sqlsqbdr.UpdatedField
		act     func() (sqlsqbdr.UpdatedField, error)
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
		{
			name:    "test5 - map string",
			isError: true,
			exp: func() sqlsqbdr.UpdatedField {
				return sqlsqbdr.UpdatedField{}
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

				h := make(map[string]Human)
				h["1"] = human
				return sqlsqbdr.BuildUpdatedField(h, sqlsqbdr.IncludeField, "name", "age")
			},
		},
		{
			name:    "test6 - interface",
			isError: true,
			exp: func() sqlsqbdr.UpdatedField {
				return sqlsqbdr.UpdatedField{}
			},
			act: func() (sqlsqbdr.UpdatedField, error) {
				var a interface{}

				return sqlsqbdr.BuildUpdatedField(a, sqlsqbdr.IncludeField, "name", "age")
			},
		},
		{
			name:    "test7 - non struct",
			isError: true,
			exp: func() sqlsqbdr.UpdatedField {
				return sqlsqbdr.UpdatedField{}
			},
			act: func() (sqlsqbdr.UpdatedField, error) {
				a := 1

				return sqlsqbdr.BuildUpdatedField(a, sqlsqbdr.IncludeField, "name", "age")
			},
		},
		{
			name: "test8 - with strip",
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
					Non    string   `json:"-" db:"-"`
				}{
					Name: "test",
					Age:  null.IntFrom(1),
					Non:  "1",
				}

				return sqlsqbdr.BuildUpdatedField(human, sqlsqbdr.IncludeField, "name", "age")
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.exp()
			act, err := tc.act()
			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, exp, act)
		})
	}
}
