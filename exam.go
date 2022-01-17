package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ExamResult struct {
	ID              int    `json:"id"`
	StudentName     string `json:"student_name"`
	StudentLastname string `json:"student_lastname"`
	ExamID          int    `json:"exam_id"`
	TeacherID       int    `json:"teacher_id"`
	Content         string `json:"content"`
	Tries           int    `json:"tries"`
}

func (e ExamResult) AddToDB() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO risultati_esami (nome_studente, cognome_studente, tentativi, ID_esame, contenuto) VALUES (?, ?, ?, ?, ?)", e.StudentName, e.StudentLastname, e.Tries, e.ExamID, e.Content)
	return err
}

type ExamToCompile struct {
	ID            int        `json:"-"`
	TeacherID     int        `json:"id_teacher"`
	Name          string     `json:"name"`
	NumOfQuestion int        `json:"num_of_question"` //[1, 50]
	ClassID       int        `json:"class_id"`
	Subject       string     `json:"subject, omitempty"`
	Difficulty    string     `json:"difficulty"` //[easy, medium, hard]
	Questions     []Question `json:"questions, omitempty"`
}

func (e *ExamToCompile) GenerateExamToCompile() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}
	var api_value string
	err = db.QueryRow("SELECT api_value FROM corsi WHERE ID=?", e.ClassID).Scan(&api_value)

	if err != nil {
		return err
	}

	resp, err := http.Get(fmt.Sprintf("https://opentdb.com/api.php?amount=%d&category=%s&difficulty=%s&type=multiple", e.NumOfQuestion, api_value, e.Difficulty))
	if err != nil {
		return err
	}

	//get the body as a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	type Response struct {
		Result []Question `json:"results"`
	}
	var res Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	e.Questions = res.Result
	return nil
}

func (e *ExamToCompile) AddToDB() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}
	result, err := db.Exec("INSERT INTO esami (nome, difficolta, ID_insegnante, numero_domande, ID_corso) VALUES (?, ?, ?, ?, ?)", e.Name, e.Difficulty, e.TeacherID, e.NumOfQuestion, e.ClassID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = int(id)
	return err
}

func (e ExamToCompile) GetExamResults(teacherID int) ([]ExamResult, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT r.ID, r.nome_studente, r.cognome_studente, r.tentativi, r.ID_esame, r.contenuto FROM risultati_esami r join esami e on ID_esame = e.ID join insegnante i on e.ID_insegnante = i.ID WHERE e.ID_insegnante = ? AND e.ID = ?;", teacherID, e.ID)
	if err != nil {
		return nil, err
	}
	var results []ExamResult
	for rows.Next() {
		var result ExamResult
		err = rows.Scan(&result.ID, &result.StudentName, &result.StudentLastname, &result.Tries, &result.ExamID, &result.Content)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

type Question struct {
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

func GetExamFromID(id int) (ExamToCompile, error) {
	var exam ExamToCompile
	db, err := ConnectToDb()
	if err != nil {
		return ExamToCompile{}, err
	}
	err = db.QueryRow("SELECT * FROM esami WHERE id = ?", id).Scan(&exam.ID, &exam.Difficulty, &exam.ClassID, &exam.NumOfQuestion, &exam.Name, &exam.TeacherID)
	return exam, err
}

func NewExam(TeacherID, numOfQuestion, classID int, name, difficulty string) (ExamToCompile, error) {
	if numOfQuestion > 50 || numOfQuestion < 1 {
		return ExamToCompile{}, fmt.Errorf("numOfQuestion must be between 1 and 50")
	}
	if difficulty != "easy" && difficulty != "medium" && difficulty != "hard" {
		return ExamToCompile{}, fmt.Errorf("difficulty must be easy, medium or hard")
	}

	exam := ExamToCompile{
		TeacherID:     TeacherID,
		NumOfQuestion: numOfQuestion,
		ClassID:       classID,
		Name:          name,
		Difficulty:    difficulty,
	}

	return exam, nil
}
