package main

import (
	"crypto/sha256"
	"fmt"
)

type Teacher struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	LastName           string `json:"lastname"`
	RegistrationNumber int    `json:"registration_number"`
	Password           string `json:"password"`
}

type Subject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func QueryTeacherByRegistrationNumber(registrationNumber int) (Teacher, error) {
	//connect to db
	db, err := ConnectToDb()
	if err != nil {
		return Teacher{}, err
	}
	var teacher Teacher
	//select every filed from teacher table where matricola = registrationNumber
	//then exec a scan over the Row returned by the query and assing the values to the teacher struct (using pointers)
	err = db.QueryRow("SELECT * FROM insegnante WHERE matricola = ?", registrationNumber).Scan(&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.RegistrationNumber, &teacher.Password)
	return teacher, err
}

func QueryTeacherByID(ID int) (Teacher, error) {
	//connect to db
	db, err := ConnectToDb()
	if err != nil {
		return Teacher{}, err
	}
	var teacher Teacher
	//select every filed from teacher table where matricola = registrationNumber
	//then exec a scan over the Row returned by the query and assing the values to the teacher struct (using pointers)
	err = db.QueryRow("SELECT * FROM insegnante WHERE ID = ?", ID).Scan(&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.RegistrationNumber, &teacher.Password)
	return teacher, err
}

func QueryAllTeachersWithAtLeastOneExam() ([]Teacher, error) {
	//connect to db
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}
	//select every filed from esami table where ID_insegnante = t.ID
	//then exec a scan over the Row returned by the query and assing the values to the exam struct
	rows, err := db.Query("SELECT i.nome, i.cognome, i.matricola FROM esami e join insegnante i on iD_insegnante = i.ID  GROUP BY iD_insegnante")
	if err != nil {
		return nil, err
	}
	var teachers []Teacher
	for rows.Next() {
		var teacher Teacher
		err = rows.Scan(&teacher.Name, &teacher.LastName, &teacher.RegistrationNumber)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func (t Teacher) AddTeacher() error {
	hashedPws := fmt.Sprintf("%x", sha256.Sum256([]byte(t.Password)))
	//connect to db
	db, err := ConnectToDb()
	if err != nil {
		return err
	}
	//insert the teacher struct into the teacher table
	_, err = db.Exec("INSERT INTO insegnante (nome, cognome, matricola, password) VALUES (?, ?, ?, ?)", t.Name, t.LastName, t.RegistrationNumber, hashedPws)
	return err
}

func (t Teacher) GetExams() ([]ExamToCompile, error) {
	//connect to db
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}
	//select every filed from esami table where ID_insegnante = t.ID
	//then exec a scan over the Row returned by the query and assing the values to the exam struct
	rows, err := db.Query("SELECT e.ID, e.ID_insegnante, e.difficolta, e.numero_domande, e.nome, e.ID_corso, c.materia FROM esami e join corsi c on ID_corso = c.ID WHERE e.id_insegnante = ?", t.ID)
	if err != nil {
		return nil, err
	}
	var exams []ExamToCompile
	for rows.Next() {
		var exam ExamToCompile
		err = rows.Scan(&exam.ID, &exam.TeacherID, &exam.Difficulty, &exam.NumOfQuestion, &exam.Name, &exam.ClassID, &exam.Subject)
		if err != nil {
			return nil, err
		}
		exams = append(exams, exam)
	}
	return exams, nil
}

func (t Teacher) GetSubjects() ([]Subject, error) {
	//connect to db
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT c.ID, c.materia FROM insegna i join corsi c on id_corso = c.ID WHERE i.id_insegnante = ?", t.ID)
	if err != nil {
		return nil, err
	}
	var subjects []Subject
	for rows.Next() {
		var subject Subject
		err = rows.Scan(&subject.ID, &subject.Name)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}

func (t Teacher) AddSubject(subject int) error {
	//connect to db
	db, err := ConnectToDb()
	if err != nil {
		return err
	}
	//insert the teacher struct into the teacher table
	_, err = db.Exec("INSERT INTO insegna (id_insegnante, id_corso) VALUES (?, ?)", t.ID, subject)
	return err
}

func GetAllSubjects() ([]Subject, error) {
	//connect to db
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT ID, materia FROM corsi")
	if err != nil {
		return nil, err
	}
	var subjects []Subject
	for rows.Next() {
		var subject Subject
		err = rows.Scan(&subject.ID, &subject.Name)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}
