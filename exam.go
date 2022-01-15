package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Exam struct {
	ID              int       `json:"id"`
	StudentName     string    `json:"student_name"`
	StudentLastname string    `json:"student_lastname"`
	ClassID         int       `json:"class_id"`
	TeacherID       int       `json:"teacher_id"`
	Date            time.Time `json:"exam_date"`
	Grade           int       `json:"grade"`
	Content         string    `json:"content"`
}

type ExamToCompile struct {
	NumOfQuestion int        `json:"num_of_question"`
	ApiCategoryID int        `json:"api_category_id"`
	Difficulty    string     `json:"difficulty"`
	Questions     []Question `json:"questions"`
}

type Question struct {
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

func (e *ExamToCompile) GenerateExamToCompile() error {
	//make a http get request to https://opentdb.com/api.php?amount=10&difficulty=easy and get the json response
	resp, err := http.Get(fmt.Sprintf("https://opentdb.com/api.php?amount=%d&category=%d&difficulty=%s&type=multiple", e.NumOfQuestion, e.ApiCategoryID, e.Difficulty))
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
