package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"restapi/internal/models"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() {
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextId++
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Luwo",
		LastName:  "Ko",
		Class:     "7B",
		Subject:   "Physics",
	}
	nextId++
}

func TeacherHandler(w http.ResponseWriter, r *http.Request) {
	// teachers/{id}
	// teachers/?key=value&querry=value2&sortby=email&sortorder=ASC
	if r.Method == http.MethodGet {
		getTeachersHandler(w, r)
	}

	if r.Method == http.MethodGet {
		addTeachersHandler(w, r)
	}

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id := strings.TrimSuffix(path, "/")

	fmt.Println("The teacher id:", id)

	querries := r.URL.Query()
	for key, value := range querries {
		fmt.Println(key, "\t|", value)
	}
	w.Write([]byte("Hello teachers route!"))
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	teacherList := make([]models.Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		teacherList = append(teacherList, teacher)
	}

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "error, teacher with id not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(teacher)
}

func addTeachersHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "error decoding json", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		newTeacher.ID = nextId
		teachers[nextId] = newTeacher
		addedTeachers[i] = newTeacher
		nextId++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(response)
}
