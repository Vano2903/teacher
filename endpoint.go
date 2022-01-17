package main

type Endpoint string

const (
	root Endpoint = "/"
	//teachers area
	TeacherPageEndpoint Endpoint = "/teacher"
	TeacherLogin        Endpoint = "/teacher/login"
	TeacherRegister     Endpoint = "/teacher/register"

	//students area
	// RequestAcceptance Endpoint = "/student/exam"
	SubmitExam        Endpoint = "/student/submit"
	AccessExam        Endpoint = "/student/access"
	RequestAccessExam Endpoint = "/student/request"

	//exam area
	GetExamResults Endpoint = "/exam/results/{id}"
	ExamPage       Endpoint = "/exam"
	AddExam        Endpoint = "/exam/add"
	ListExams      Endpoint = "/exam/list"
)

func (e Endpoint) String() string {
	return string(e)
}
