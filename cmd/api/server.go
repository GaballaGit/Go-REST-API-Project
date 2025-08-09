package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	mw "restapi/internal/api/middlewares"
)

type Teacher struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Class     string `json:"class"`
	Subject   string `json:"subject"`
}

var (
	teachers = make(map[int]Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() {
	teachers[nextId] = Teacher{
		ID:        nextId,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextId++
	teachers[nextId] = Teacher{
		ID:        nextId,
		FirstName: "Luwo",
		LastName:  "Ko",
		Class:     "7B",
		Subject:   "Physics",
	}
	nextId++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	teacherList := make([]Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		teacherList = append(teacherList, teacher)
	}

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []Teacher `json:"data"`
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

	var newTeachers []Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "error decoding json", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		newTeacher.ID = nextId
		teachers[nextId] = newTeacher
		addedTeachers[i] = newTeacher
		nextId++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(response)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQUEST:", r.Method)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Hello root server!"))
	/*fmt.Println("form:", r.Form)
	process := make(map[string]interface{})
	for key, value := range r.Form {
		process[key] = value
	}
	fmt.Println("Processed map:", process)
	*/
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
	}
	fmt.Println("body:", string(body))
	fmt.Println("context:", r.Context())
	fmt.Println("content len:", r.ContentLength)
	fmt.Println(r)
}

func handleTeachers(w http.ResponseWriter, r *http.Request) {
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

func handleStudents(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello students route!"))
}

func handleExecs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Execs route!"))
}

func main() {

	PORT := 8080

	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := mw.NewRateLimiter(6, time.Minute)

	hpp := mw.HPPOptions{
		CheckQuerry:             true,
		CheckBody:               true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		Whitelist:               []string{"sortby", "sortorder", "name", "age", "class"},
	}

	//secureMux := mw.Cors(rl.MiddleWare(mw.ResponseTime(mw.Compression(mw.SecurityHeader(mw.Hpp(hpp)(mux))))))
	secureMux := applyMiddlewares(
		mux,
		mw.Hpp(hpp),
		mw.Compression,
		mw.SecurityHeader,
		mw.ResponseTime,
		rl.MiddleWare,
		mw.Cors,
	)

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", PORT),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/teachers/", handleTeachers)
	mux.HandleFunc("/students/", handleStudents)
	mux.HandleFunc("/execs/", handleExecs)

	fmt.Println("server started at port:", PORT)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalf("error starting server: %s", err)
		return
	}
}

// Make middlewares cleaner
type MiddleWare func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...MiddleWare) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
