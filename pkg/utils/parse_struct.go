package utils

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"restapi/internal/models"
)

func GetStructValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()
	var values []interface{}
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" {
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	log.Println("Values:", values)
	return values
}

// Takes a model and gets the current db value. Then iterate over update
// map to upadate the model. Call sqlconnect.GenerateUpdateQuery() after this.
func PatchTeacherModel(db *sql.DB, tabel string, id int, model *models.Teacher, update map[string]interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %d", tabel, id)
	err := db.QueryRow(query).Scan(
		&model.ID,
		&model.FirstName,
		&model.LastName,
		&model.Email,
		&model.Class,
		&model.Subject,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorHandler(err, "teacher not found")
		}
		return ErrorHandler(err, "error retrieving teacher")
	}

	modelVal := reflect.ValueOf(model).Elem()
	modelTyp := modelVal.Type()

	for k, v := range update {
		if k == "id" {
			continue
		}
		for i := 0; i < modelVal.NumField(); i++ {
			field := modelTyp.Field(i)

			if field.Tag.Get("json") == k+",omitempty" {
				fieldVal := modelVal.Field(i)

				if fieldVal.CanSet() {
					val := reflect.ValueOf(v)

					if val.Type().ConvertibleTo(fieldVal.Type()) {
						fieldVal.Set(val.Convert(fieldVal.Type()))
					} else {
						msg := fmt.Sprintf("cannot convert %v to %v", val.Type(), fieldVal.Type())
						return ErrorHandler(fmt.Errorf("%s", msg), msg)
					}
				}
				break
			}
		}
	}
	return nil
}

// Takes a model and gets the current db value. Then iterate over update
// map to upadate the model. Call sqlconnect.GenerateUpdateQuery() after this.
func PatchStudentModel(db *sql.DB, tabel string, id int, model *models.Student, update map[string]interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %d", tabel, id)
	err := db.QueryRow(query).Scan(
		&model.ID,
		&model.FirstName,
		&model.LastName,
		&model.Email,
		&model.Class,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorHandler(err, "teacher not found")
		}
		return ErrorHandler(err, "error retrieving teacher")
	}

	modelVal := reflect.ValueOf(model).Elem()
	modelTyp := modelVal.Type()

	for k, v := range update {
		if k == "id" {
			continue
		}
		for i := 0; i < modelVal.NumField(); i++ {
			field := modelTyp.Field(i)

			if field.Tag.Get("json") == k+",omitempty" {
				fieldVal := modelVal.Field(i)

				if fieldVal.CanSet() {
					val := reflect.ValueOf(v)

					if val.Type().ConvertibleTo(fieldVal.Type()) {
						fieldVal.Set(val.Convert(fieldVal.Type()))
					} else {
						msg := fmt.Sprintf("cannot convert %v to %v", val.Type(), fieldVal.Type())
						return ErrorHandler(fmt.Errorf("%s", msg), msg)
					}
				}
				break
			}
		}
	}
	return nil
}
