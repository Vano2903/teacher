package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type Post struct {
	Name               string `json:"name, omitempty"`
	LastName           string `json:"lastname, omitempty"`
	RegistrationNumber int    `json:"registration_number, omitempty"`
	Password           string `json:"password, omitempty"`
	ExamID             int    `json:"exam_id, omitempty"`
}

//middlewars
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var (
			found bool
		)
		//read value of cookie called jwt
		for _, cookie := range r.Cookies() {
			if cookie.Name == "JWT" {
				found = true
				break
			}
		}

		if !found {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "missing 'JWT' cookie")
			return
		}

		next.ServeHTTP(w, r)
	})
}

//teacher's handlers
func TeacherPage(w http.ResponseWriter, r *http.Request) {
	var jwt string

	//read value of cookie called jwt
	for _, cookie := range r.Cookies() {
		if cookie.Name == "JWT" {
			jwt = cookie.Value
			break
		}
	}

	teacherJWT, err := ParseToken(jwt)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Invalid authorization token: "+err.Error())
		return
	}

	data := struct {
		Name     string
		LastName string
	}{
		Name:     teacherJWT.Name,
		LastName: teacherJWT.Lastname,
	}

	tmpl, err := template.ParseFiles("pages/teacher.html")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusServiceUnavailable, "Internal server error: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func LoginTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//read from the post body the json data and fill the post struct
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusBadRequest, "Invalid json")
	}
	//hash the password with sha256
	hashedPassword := sha256.Sum256([]byte(post.Password))
	teacher, err := QueryTeacherByRegistrationNumber(post.RegistrationNumber)
	if err != nil {
		//internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}
	if teacher.Password == fmt.Sprintf("%x", hashedPassword) {
		//return the teacher struct as json
		// json.NewEncoder(w).Encode(teacher)
		teacherByte, err := json.Marshal(teacher)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
			return
		}

		Unsignedtoken := CustomClaims{
			RegistrationNumber: teacher.RegistrationNumber,
			Lastname:           teacher.LastName,
			Name:               teacher.Name,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
				Issuer:    "vano-jwt-teachers",
			},
		}

		token, err := NewSignedToken(Unsignedtoken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
			return
		}
		// w.Header().Add("Authorization", "Bearer "+token)

		cookie := &http.Cookie{
			Name:     "JWT",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * time.Duration(1)),
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
		w.Header().Add("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "teacher":%s, "error": false}`, http.StatusAccepted, "Successfully logged in", string(teacherByte))

		log.Printf("successful login by %s %s %d \n", teacher.Name, teacher.LastName, teacher.RegistrationNumber) //logging
		return
	}
	//return unauthorized
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Incorrect credentials")
}

func RegisterTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//read from the post body the json data and fill the post struct
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusBadRequest, "Invalid json")
	}

	//fill teacher's struct
	var t Teacher
	t.Name = post.Name
	t.LastName = post.LastName
	t.RegistrationNumber = post.RegistrationNumber
	t.Password = post.Password
	//check if the registration number is already in use
	_, err = QueryTeacherByRegistrationNumber(post.RegistrationNumber)
	if err == nil {
		//internal server error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusBadRequest, "a teacher with this registration number already exists")
		return
	}

	//insert the teacher struct into the teacher table
	err = t.AddTeacher()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}
	//return the teacher struct as jso
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": false}`, http.StatusCreated, "Teacher successfully registered, you can do the login now")
	log.Printf("successful registered by %s %s %d \n", post.Name, post.LastName, post.RegistrationNumber) //logging
}

//student's handlers
func RequestAccessExamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusBadRequest, "Invalid json")
		return
	}

	UnsignedJWT := NewCustomClaims(p.Name, p.LastName, p.ExamID, time.Now().Add(time.Hour*time.Duration(1)).Unix())
	token, err := NewSignedToken(UnsignedJWT)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}
	// w.Header().Add("Authorization", "Bearer "+token)

	cookie := &http.Cookie{
		Name:     "JWT",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.Header().Add("Authorization", "Bearer "+token)
	fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": false}`, http.StatusOK, "Successfully requested access to exam, you have 1 hour to complete the test from now, you can access only one time to the resource but you can do as many tries as you wish (by staying in the time limit)")
}

func AccessExamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var jwt string

	//read value of cookie called jwt
	for _, cookie := range r.Cookies() {
		if cookie.Name == "JWT" {
			jwt = cookie.Value
			break
		}
	}

	//convert the jwt
	value, err := ParseToken(jwt)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Invalid authorization token: "+err.Error())
		return
	}

	//check if the user already accessed the exam
	if value.Accessed {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Access already happened, must ask for a new access token")
		return
	}

	exam, err := GetExamFromID(value.ExamID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	err = exam.GenerateExamToCompile()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	questions, err := json.Marshal(exam.Questions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	value.Accessed = true
	token, err := NewSignedToken(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     "JWT",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * time.Duration(1)),
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.Header().Add("Authorization", "Bearer "+token)

	fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": false, "questions":%s}`, http.StatusOK, "access to exam granted", questions)
}

func SubmitExamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var jwt string

	//read value of cookie called jwt
	for _, cookie := range r.Cookies() {
		if cookie.Name == "JWT" {
			jwt = cookie.Value
			break
		}
	}

	//convert the jwt
	value, err := ParseToken(jwt)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Invalid authorization token: "+err.Error())
		return
	}

	exam, err := GetExamFromID(value.ExamID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	type ExamPost struct {
		Content string `json:"content"`
		Tries   int    `json:"tries"`
	}

	var post ExamPost
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusBadRequest, "Invalid json")
		return
	}

	result := ExamResult{
		ExamID:          value.ExamID,
		StudentName:     value.Name,
		StudentLastname: value.Lastname,
		TeacherID:       exam.TeacherID,
		Content:         post.Content,
		Tries:           post.Tries,
	}

	if value.Submitted {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusBadRequest, "You have already submitted the exam")
		return
	}

	if err := result.AddToDB(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	value.Submitted = true
	token, err := NewSignedToken(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     "JWT",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.Header().Add("Authorization", "Bearer "+token)
	fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": false}`, http.StatusOK, "exam submitted")
}

//exam's handlers
func AddExamHandler(w http.ResponseWriter, r *http.Request) {
	var jwt string

	//read value of cookie called jwt
	for _, cookie := range r.Cookies() {
		if cookie.Name == "JWT" {
			jwt = cookie.Value
			break
		}
	}

	teacherJWT, err := ParseToken(jwt)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Invalid authorization token: "+err.Error())
		return
	}

	var post ExamToCompile
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusBadRequest, "Invalid json")
		return
	}

	teacher, err := QueryTeacherByRegistrationNumber(teacherJWT.RegistrationNumber)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Invalid authorization token: "+err.Error())
		return
	}

	exam, err := NewExam(teacher.ID, post.NumOfQuestion, post.ClassID, post.Name, post.Difficulty)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	if err := exam.AddToDB(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": false}`, http.StatusCreated, "Exam successfully added, the ID is: "+strconv.Itoa(exam.ID))
}

func GetResultsOfExamHandler(w http.ResponseWriter, r *http.Request) {
	var jwt string

	//read value of cookie called jwt
	for _, cookie := range r.Cookies() {
		if cookie.Name == "JWT" {
			jwt = cookie.Value
			break
		}
	}

	//parse the jwt
	teacherJWT, err := ParseToken(jwt)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Invalid authorization token: "+err.Error())
		return
	}

	//get the id from get
	examID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusBadRequest, "Invalid exam ID")
		return
	}

	//get the exam from the id given in the url
	exam, err := GetExamFromID(examID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	//get the teacher using the value in the jwt
	teacher, err := QueryTeacherByRegistrationNumber(teacherJWT.RegistrationNumber)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Invalid authorization token: "+err.Error())
		return
	}

	//check if the teacher is the creator of the exam
	//and get the results of the exam
	//if the examid or the teacherid are invalid it will return a emtpy slice
	results, err := exam.GetExamResults(teacher.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	resultsJSON, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": false, "results":%s}`, http.StatusOK, "Students successfully retrieved", resultsJSON)
}

func ListExamsHandler(w http.ResponseWriter, r *http.Request) {
	var jwt string

	//read value of cookie called jwt
	for _, cookie := range r.Cookies() {
		if cookie.Name == "JWT" {
			jwt = cookie.Value
			break
		}
	}

	//parse the jwt
	teacherJWT, err := ParseToken(jwt)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Invalid authorization token: "+err.Error())
		return
	}

	//get the teacher using the value in the jwt
	teacher, err := QueryTeacherByRegistrationNumber(teacherJWT.RegistrationNumber)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusUnauthorized, "Invalid authorization token: "+err.Error())
		return
	}

	//get the exams of the teacher
	//if the teacherid is invalid it will return a emtpy slice
	exams, err := teacher.GetExams()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": true}`, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	examsJSON, _ := json.Marshal(exams)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"code": %d, "msg":"%s", "error": false, "exams":%s}`, http.StatusOK, "Exams successfully retrieved", examsJSON)
}

func main() {
	r := mux.NewRouter()
	//teacher handlers
	r.Handle(TeacherPageEndpoint.String(), JWTAuthMiddleware(http.HandlerFunc(TeacherPage))).Methods("GET")
	r.HandleFunc(TeacherLogin.String(), LoginTeacherHandler).Methods("POST")
	r.HandleFunc(TeacherRegister.String(), RegisterTeacherHandler).Methods("POST")

	//exam handlers
	r.Handle(AddExam.String(), JWTAuthMiddleware(http.HandlerFunc(AddExamHandler))).Methods("POST")
	r.Handle(GetExamResults.String(), JWTAuthMiddleware(http.HandlerFunc(GetResultsOfExamHandler))).Methods("GET")
	r.Handle(ListExams.String(), JWTAuthMiddleware(http.HandlerFunc(ListExamsHandler))).Methods("GET")

	//student handlers
	r.HandleFunc(RequestAccessExam.String(), RequestAccessExamHandler).Methods("POST")
	r.Handle(AccessExam.String(), JWTAuthMiddleware(http.HandlerFunc(AccessExamHandler))).Methods("GET")
	r.Handle(SubmitExam.String(), JWTAuthMiddleware(http.HandlerFunc(SubmitExamHandler))).Methods("POST")
	http.ListenAndServe(":8080", r)
}
