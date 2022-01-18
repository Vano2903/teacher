package main

type Endpoint string

const (
	AllSubjects Endpoint = "/subjects"

	//teachers area
	TeacherPageEndpoint Endpoint = "/teacher"
	TeacherLogin        Endpoint = "/teacher/login"
	TeacherRegister     Endpoint = "/teacher/register"
	GetSubjects         Endpoint = "/teacher/subjects"
	AddSubject          Endpoint = "/teacher/subjects/add"
	//students area
	// RequestAcceptance Endpoint = "/student/exam"
	RequestTeacherExam Endpoint = "/student/exam/teachers"
	RequestExam        Endpoint = "/student/exam/{teacher_id}/exams"
	GetExamPage        Endpoint = "/student/exam"
	SubmitExam         Endpoint = "/student/submit"
	AccessExam         Endpoint = "/student/access"
	RequestAccessExam  Endpoint = "/student/request"

	//exam area
	GetExamResults Endpoint = "/exam/results/{id}"
	ExamPage       Endpoint = "/exam"
	AddExam        Endpoint = "/exam/add"
	ListExams      Endpoint = "/exam/list/{teacher_id}"
	GetAllTeachers Endpoint = "/exam/teachers"
)

func (e Endpoint) String() string {
	return string(e)
}
