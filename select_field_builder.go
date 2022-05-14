package sqlsqbdr

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func BuildSelectFields(entity interface{}, fields ...string) (string, error) {
	typeOf := reflect.TypeOf(entity)
	valueOf := reflect.ValueOf(entity)

	if entity == nil {
		return "", fmt.Errorf("empty")
	}
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}

	if typeOf.Kind() != reflect.Struct {
		return "", errors.New("not a struct")
	}

	selectFields := []string{}

	for _, field := range fields {
		for i := 0; i < valueOf.NumField(); i++ {
			f := typeOf.Field(i)
			name := strings.Split(f.Tag.Get("db"), ",")[0]

			if field == name {
				selectFields = append(selectFields, name)
			}
		}
	}

	return strings.Join(selectFields, ","), nil
}
