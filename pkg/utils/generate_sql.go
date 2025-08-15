package utils

import (
	"fmt"

	"reflect"

	"strings"
)

// Generates an insert querry given a table and new model
//
// INSET INTO table (first_name, last_name,...) VALUE (?, ?, ?, ...)
func GenerateInsertQuery(table string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")

		if dbTag != "" && dbTag != "id" {
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			placeholders += "?"
		}
	}

	fmt.Println(columns, "\n--\n", placeholders)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, columns, placeholders)
}

// Generate an insert query
//
// UPDATE table SET age = ?, name = ?,... WHERE id = ?
func GenerateUpdateQuery(table string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	modelVal := reflect.ValueOf(model)
	updates := ""

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		fmt.Println(i, dbTag)

		if dbTag != "" && dbTag != "id" {
			if updates != "" {
				updates += ", "
			}
			updates += fmt.Sprintf("%s = \"%s\"", dbTag, modelVal.Field(i).Interface())
		}
	}

	id := modelVal.Field(0).Int()

	fmt.Printf("UPDATE %s SET %s WHERE id = %d\n", table, updates, id)
	return fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", table, updates, id)
}

// Generate a delete sql query
//
// DELETE FROM table WHERE id = ?
func GenerateDeleteQuery(table string, model interface{}) string {
	modelVal := reflect.ValueOf(model)
	//modelTyp := reflect.TypeOf(model)
	id := modelVal.Field(0).Int()

	return fmt.Sprintf("DELETE FROM %s WHERE id = %d", table, id)
}
