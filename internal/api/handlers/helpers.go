package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"restapi/pkg/utils"
	"strings"
)

// Check if there exists a blank field. Returns and error if so.
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

// Return a string slice of field names.
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

// Add sorting to the query
func sortBy(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			fmt.Println(len(parts), parts)
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			fmt.Println(isValidSortField(field), isValidSortOrder(order))
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			fmt.Println(field, order)
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}
	}
	return query
}

// Return if sort order is valid ("asc" or "desc")
func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}
