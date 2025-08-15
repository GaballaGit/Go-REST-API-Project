package handlers

import (
	"fmt"
	"reflect"
	"restapi/pkg/utils"
	"strings"
)

func checkBlankFields(model interface{}) error {
	val := reflect.ValueOf(model)
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Name == "id" {
			continue
		}
		fieldVal := val.Field(i)
		if fieldVal.Kind() == reflect.String && fieldVal.String() == "" {
			return utils.ErrorHandler(fmt.Errorf("invalid field in models"), "all fields are required")
		}
	}
	return nil
}

func getFieldNames(model interface{}) []string {
	modVal := reflect.ValueOf(model)
	modTyp := modVal.Type()
	fields := []string{}

	for i := 0; i < modTyp.NumField(); i++ {
		field := modTyp.Field(i)
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldToAdd)
	}
	return fields
}
