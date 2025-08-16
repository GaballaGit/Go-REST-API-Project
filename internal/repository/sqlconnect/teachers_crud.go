package sqlconnect

import (
	"database/sql"
	"fmt"
	"net/http"
	"restapi/internal/models"
	"restapi/pkg/utils"
)

func DeleteTeachersFromDb(w http.ResponseWriter, ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error opening database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error preparing transaction")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "error preparing delete statment")
	}
	defer stmt.Close()

	var deletedIds []int
	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "error executing statement")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "error retrieving delete results")
		}

		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error committing transaction")
	}

	if len(deletedIds) < 1 {
		return nil, utils.ErrorHandler(err, "error, ids do not exist")
	}

	return deletedIds, nil
}

func DeleteOneTeacher(w http.ResponseWriter, id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "error opening up db")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "error deleting request from database")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "error getting rows affected")
	}

	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "teacher was not found")
	}
	return nil
}

func PatchOneTeacher(w http.ResponseWriter, id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error opening up database")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = utils.PatchTeacherModel(db, "teachers", id, &existingTeacher, updates)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error patching model")
	}

	update := utils.GenerateUpdateQuery("teachers", existingTeacher)
	_, err = db.Exec(update)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error updating teacher")
	}
	return existingTeacher, nil
}

// Takes a map of fields to be patched
func PatchTeachers(w http.ResponseWriter, updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "unable to connect to database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "error starting transaction")
	}

	fmt.Println(updates)
	for _, update := range updates {
		idFloat, ok := update["id"].(float64)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "invalid teacher id")
		}

		id := int(idFloat)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "error converting id to int")
		}

		var teacherFromDb models.Teacher
		err = utils.PatchTeacherModel(db, "teachers", id, &teacherFromDb, update)
		if err != nil {
			return utils.ErrorHandler(err, "error updating teacher struct")
		}

		update := utils.GenerateUpdateQuery("teachers", teacherFromDb)
		_, err = tx.Exec(update)
		if err != nil {
			fmt.Println(teacherFromDb)
			tx.Rollback()
			return utils.ErrorHandler(err, "error updating teacher")
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "error comitting transaction")
	}
	return nil
}

func UpdateTeacher(w http.ResponseWriter, id int, updateTeacher models.Teacher) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "error opening up database")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject,
	)
	if err == sql.ErrNoRows {
		return utils.ErrorHandler(err, "teacher not found")
	} else if err != nil {
		return utils.ErrorHandler(err, "error retrieving teacher from database")
	}

	updateTeacher.ID = existingTeacher.ID
	updateQuery := utils.GenerateUpdateQuery("teachers", updateTeacher)

	_, err = db.Exec(updateQuery)
	if err != nil {
		return utils.ErrorHandler(err, "error updating teacher")
	}
	return nil
}

func AddTeacher(w http.ResponseWriter, newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "error opening up database", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "error opening up database")
	}
	defer db.Close()

	//stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("teachers", models.Teacher{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "error in preparing SQL query")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error inserting data into database")
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "error getting last insert id")
		}
		newTeacher.ID = int(lastId)
		addedTeachers[i] = newTeacher
	}
	return addedTeachers, nil
}

func GetOneTeacher(w http.ResponseWriter, id int) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error trying to open SQL database")
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		return models.Teacher{}, utils.ErrorHandler(err, "error teachers not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error getting teacher from database")
	}
	return teacher, nil
}
