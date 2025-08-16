package sqlconnect

import (
	"database/sql"
	"fmt"
	"net/http"
	"restapi/internal/models"
	"restapi/pkg/utils"
)

func GetOneStudent(w http.ResponseWriter, id int) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "error trying to open SQL database")
	}
	defer db.Close()

	var student models.Student
	err = db.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		return models.Student{}, utils.ErrorHandler(err, "error student not found")
	} else if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "error getting student from database")
	}
	return student, nil
}

func AddStudent(w http.ResponseWriter, newStudent []models.Student) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "error opening up database", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "error opening up database")
	}
	defer db.Close()

	//stmt, err := db.Prepare("INSERT INTO students (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("students", models.Student{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "error in preparing SQL query")
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, len(newStudent))
	for i, addStudent := range newStudent {
		fmt.Println(addStudent.FirstName, addStudent.LastName, addStudent.Email, addStudent.Class)
		res, err := stmt.Exec(addStudent.FirstName, addStudent.LastName, addStudent.Email, addStudent.Class)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error inserting data into database")
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "error getting last insert id")
		}
		addStudent.ID = int(lastId)
		addedStudents[i] = addStudent
	}
	return addedStudents, nil
}

func UpdateStudent(w http.ResponseWriter, id int, updateStudent models.Student) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "error opening up database")
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(
		&existingStudent.ID,
		&existingStudent.FirstName,
		&existingStudent.LastName,
		&existingStudent.Email,
		&existingStudent.Class,
	)
	if err == sql.ErrNoRows {
		return utils.ErrorHandler(err, "student not found")
	} else if err != nil {
		return utils.ErrorHandler(err, "error retrieving student from database")
	}

	updateStudent.ID = existingStudent.ID
	updateQuery := utils.GenerateUpdateQuery("students", updateStudent)

	_, err = db.Exec(updateQuery)
	if err != nil {
		return utils.ErrorHandler(err, "error updating student")
	}
	return nil
}

func PatchStudent(w http.ResponseWriter, updates []map[string]interface{}) error {
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
			return utils.ErrorHandler(err, "invalid student id")
		}

		id := int(idFloat)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "error converting id to int")
		}

		var studentFromDb models.Student
		err = utils.PatchStudentModel(db, "students", id, &studentFromDb, update)
		if err != nil {
			return utils.ErrorHandler(err, "error updating student struct")
		}

		update := utils.GenerateUpdateQuery("students", studentFromDb)
		_, err = tx.Exec(update)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "error updating student")
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "error comitting transaction")
	}
	return nil
}

func PatchOneStudent(w http.ResponseWriter, id int, updates map[string]interface{}) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "error opening up database")
	}
	defer db.Close()

	var existingStudent models.Student
	err = utils.PatchStudentModel(db, "students", id, &existingStudent, updates)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "error patching model")
	}

	update := utils.GenerateUpdateQuery("students", existingStudent)
	_, err = db.Exec(update)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "error updating student")
	}
	return existingStudent, nil
}

func DeleteOneStudent(w http.ResponseWriter, id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "error opening up db")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "error deleting request from database")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "error getting rows affected")
	}

	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "student was not found")
	}
	return nil
}

func DeleteStudentsFromDb(w http.ResponseWriter, ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error opening database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error preparing transaction")
	}

	stmt, err := tx.Prepare("DELETE FROM students WHERE id = ?")
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
