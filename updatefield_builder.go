package sqlsqbdr

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type UpdatedField struct {
	Name  []string
	Value []interface{}
}

func BuildUpdatedField(entity interface{}, fieldSelectorType TypeFieldSelect, fields ...string) (UpdatedField, error) {
	fieldMap := MultipleStringToMap(fields)
	typeOf := reflect.TypeOf(entity)
	valueOf := reflect.ValueOf(entity)
	if typeOf.Kind() != reflect.Struct {
		return UpdatedField{}, errors.New("not a struct")
	}

	var (
		fieldName  []string
		fieldValue []interface{}
	)

	for i := 0; i < valueOf.NumField(); i++ {
		f := typeOf.Field(i)
		name := strings.Split(f.Tag.Get("db"), ",")[0]
		value := valueOf.Field(i)
		switch fieldSelectorType {
		case ExcludeField:
			if fieldMap[name] == name || len(fieldMap) == 0 {
				fieldName = append(fieldName, fmt.Sprintf("%s = ?", name))
				fieldValue = append(fieldValue, value.Interface())
			}
		default:
			if fieldMap[name] == name || len(fieldMap) == 0 {
				fieldName = append(fieldName, fmt.Sprintf("%s = ?", name))
				fieldValue = append(fieldValue, value.Interface())
			}
		}
	}

	return UpdatedField{
		Name:  fieldName,
		Value: fieldValue,
	}, nil
}
