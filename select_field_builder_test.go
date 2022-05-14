package sqlsqbdr_test

import (
	"testing"

	"github.com/nurcahyaari/sqlsqbdr"
	"github.com/stretchr/testify/assert"
)

func TestBuildSelectFields(t *testing.T) {

	testCase := []struct {
		name    string
		isError bool
		exp     func() string
		act     func() (string, error)
	}{
		{
			name: "Test1",
			exp: func() string {
				return "name,age"
			},
			act: func() (string, error) {
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
				return sqlsqbdr.BuildSelectFields(human, "name", "age")
			},
		},
		{
			name: "Test2 - struct",
			exp: func() string {
				return "name,age"
			},
			act: func() (string, error) {
				return sqlsqbdr.BuildSelectFields(struct {
					Name    string `json:"name" db:"name"`
					Age     int    `json:"age" db:"age"`
					Gender  string `json:"gender" db:"gender"`
					Address string `json:"address" db:"address"`
				}{
					Name:    "test",
					Age:     0,
					Address: "Test1",
				}, "name", "age")
			},
		},
		{
			name: "Test3",
			exp: func() string {
				return "name,age"
			},
			act: func() (string, error) {
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
				return sqlsqbdr.BuildSelectFields(human, "name", "age", "test")
			},
		},
		{
			name:    "Test4 - non struct",
			isError: true,
			exp: func() string {
				return ""
			},
			act: func() (string, error) {
				return sqlsqbdr.BuildSelectFields(1, "name", "age", "test")
			},
		},
		{
			name:    "Test5 - non struct interface",
			isError: true,
			exp: func() string {
				return ""
			},
			act: func() (string, error) {
				var a interface{}
				return sqlsqbdr.BuildSelectFields(a, "name", "age", "test")
			},
		},
		{
			name:    "Test6 - non struct interface",
			isError: true,
			exp: func() string {
				return ""
			},
			act: func() (string, error) {
				var a interface{}
				a = 1
				return sqlsqbdr.BuildSelectFields(a, "name", "age", "test")
			},
		},
		{
			name:    "Test6 - non struct interface",
			isError: true,
			exp: func() string {
				return ""
			},
			act: func() (string, error) {
				var a map[string]interface{}
				return sqlsqbdr.BuildSelectFields(a, "name", "age", "test")
			},
		},
		{
			name:    "Test7 - map struct",
			isError: true,
			exp: func() string {
				return ""
			},
			act: func() (string, error) {
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
				return sqlsqbdr.BuildSelectFields(h, "name", "age", "test")
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
