package main

type Endpoint string

const (
	root Endpoint = "/"
	//teachers area
	teacherLogin    Endpoint = "/teacher/login"
	teacherRegister Endpoint = "/teacher/register"

	//students area
	requestAcceptance Endpoint = "/student/exam"

	//exam area
	examList Endpoint = "/exam/list"
)
