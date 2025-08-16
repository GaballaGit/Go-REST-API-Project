package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"strconv"
)

func StudentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("Get students"))
	case "POST":
		w.Write([]byte("Post students"))
	case "PUT":
		w.Write([]byte("Put students"))
	case "DELETE":
		w.Write([]byte("Delete students"))
	}
}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "unable to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "SELECT * FROM students WHERE 1=1"
	var args []interface{}

	query, args = addStudentFilter(r, query, args)
	query = sortBy(r, query)

	// get rows
	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "error querying db", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	studentsList := make([]models.Student, 0)

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			http.Error(w, "error scanning database results", http.StatusInternalServerError)
			return
		}
		studentsList = append(studentsList, student)
	}
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(studentsList),
		Data:   studentsList,
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	student, err := sqlconnect.GetOneStudent(w, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(student)
}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	var newStudent []models.Student
	var rawStudent []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &rawStudent)
	if err != nil {
		http.Error(w, "error decoding json", http.StatusBadRequest)
		return
	}

	fields := getFieldNames(models.Student{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	fmt.Println(rawStudent)
	for _, teacher := range rawStudent {
		for key := range teacher {
			if _, ok := allowedFields[key]; !ok {
				http.Error(w, "unnaccepable field found in request", http.StatusBadRequest)
				return
			}
		}
	}

	err = json.Unmarshal(body, &newStudent)
	if err != nil {
		http.Error(w, "error decoding json", http.StatusBadRequest)
		return
	}

	for _, teacher := range newStudent {
		err = checkBlankFields(teacher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	addedStudent, err := sqlconnect.AddStudent(w, newStudent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudent),
		Data:   addedStudent,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateStudentsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student Id", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}

	var updateStudent models.Student
	err = json.Unmarshal(body, &updateStudent)
	if err != nil {
		http.Error(w, "error decoding request body", http.StatusInternalServerError)
		return
	}

	err = sqlconnect.UpdateStudent(w, id, updateStudent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateStudent)
}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchStudent(w, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PatchOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student Id", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	existingStudent, err := sqlconnect.PatchOneStudent(w, id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingStudent)
}

func DeleteOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid student request", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneStudent(w, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Status 204, with no response body
	//w.WriteHeader(http.StatusNoContent)

	// With response body
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Student succesfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "error retrieving body values", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteStudentsFromDb(w, ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		Deletedids []int  `json:"deleted_ids"`
	}{
		Status:     "Students succesfully deleted",
		Deletedids: deletedIds,
	}

	json.NewEncoder(w).Encode(response)
}
func addStudentFilter(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}
