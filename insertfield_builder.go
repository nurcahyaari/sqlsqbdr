package sqlsqbdr

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type InsertField struct {
	Name        []string
	Placeholder []string
	Values      []interface{}
}

func createInsertField(typeOf reflect.Type, valueOf reflect.Value, fieldSelectorType TypeFieldSelect, fieldMap map[string]string) ([]string, []string, []interface{}) {
	var fieldName []string
	var fieldValue []interface{}
	var fieldPlaceholder []string
	for i := 0; i < valueOf.NumField(); i++ {
		f := typeOf.Field(i)
		name := strings.Split(f.Tag.Get("db"), ",")[0]
		value := valueOf.Field(i)
		switch fieldSelectorType {
		case ExcludeField:
			if (fieldMap[name] != name || len(fieldMap) == 0) && name != "-" {
				fieldName = append(fieldName, name)
				fieldPlaceholder = append(fieldPlaceholder, "?")
				fieldValue = append(fieldValue, value.Interface())
			}
		default:
			if (fieldMap[name] == name || len(fieldMap) == 0) && name != "-" {
				fieldName = append(fieldName, name)
				fieldPlaceholder = append(fieldPlaceholder, "?")
				fieldValue = append(fieldValue, value.Interface())
			}
		}
	}

	return fieldName, fieldPlaceholder, fieldValue
}

func BuildInsertField(entities interface{}, fieldSelectorType TypeFieldSelect, fields ...string) (InsertField, error) {
	fieldMap := MultipleStringToMap(fields)

	typeOfEntities := reflect.TypeOf(entities)
	valuesOfEntities := reflect.ValueOf(entities)

	var (
		fieldName         []string
		fieldPlaceholders []string
		fieldValues       []interface{}
	)

	if typeOfEntities.Kind() == reflect.Slice || typeOfEntities.Kind() == reflect.Array {
		for i := 0; i < valuesOfEntities.Len(); i++ {
			valueOfEntity := valuesOfEntities.Index(i).Interface()
			typeOf := reflect.TypeOf(valueOfEntity)
			valueOf := reflect.ValueOf(valueOfEntity)

			if typeOf.Kind() == reflect.Ptr {
				typeOf = typeOf.Elem()
				valueOf = valueOf.Elem()
			}

			if typeOf.Kind() != reflect.Struct {
				return InsertField{}, errors.New("not a struct")
			}

			fname, fplaceholder, fvalue := createInsertField(typeOf, valueOf, fieldSelectorType, fieldMap)
			if fname == nil || fplaceholder == nil || fvalue == nil {
				continue
			}

			fieldName = fname
			fieldPlaceholders = append(fieldPlaceholders, fmt.Sprintf("(%s)", strings.Join(fplaceholder, ",")))
			fieldValues = append(fieldValues, fvalue...)
		}
	} else {
		if typeOfEntities.Kind() == reflect.Ptr {
			typeOfEntities = typeOfEntities.Elem()
			valuesOfEntities = valuesOfEntities.Elem()
		}

		if typeOfEntities.Kind() != reflect.Struct {
			return InsertField{}, errors.New("not a struct")
		}
		fname, fplaceholder, fvalue := createInsertField(typeOfEntities, valuesOfEntities, fieldSelectorType, fieldMap)
		if fname == nil || fplaceholder == nil || fvalue == nil {
			return InsertField{}, nil
		}

		fieldName = fname
		fieldPlaceholders = append(fieldPlaceholders, fmt.Sprintf("(%s)", strings.Join(fplaceholder, ",")))
		fieldValues = append(fieldValues, fvalue...)
	}

	return InsertField{
		Name:        fieldName,
		Placeholder: fieldPlaceholders,
		Values:      fieldValues,
	}, nil
}
